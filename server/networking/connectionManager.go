package networking

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/Shiwang0-0/multiplayertetris/protocol"
)

type Client struct {
	ID      int
	conn    net.Conn
	outChan chan protocol.Msg
	done    chan struct{} // closed once, only by the owning reader goroutine (dont belive on other go routines to close, might write to a closed channel because of race condition)
}

type ConnectionManager struct {
	clients map[int]Client
	nextId  int
	mu      sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: map[int]Client{},
	}
}

// in ConnectionManager, a new setup method:
func (cm *ConnectionManager) Accept(conn net.Conn) Client {
	clientId := cm.getNextId()
	client := Client{
		ID:      clientId,
		conn:    conn,
		outChan: make(chan protocol.Msg, 16),
		done:    make(chan struct{}),
	} // buffered, see #6
	cm.mu.Lock()
	cm.clients[clientId] = client
	cm.mu.Unlock()
	return client
}

// Reading the client messages
func (cm *ConnectionManager) HandleConnectionRead(client Client, rm *RoomManager) {
	defer func() {
		client.conn.Close()
		cm.mu.Lock()
		delete(cm.clients, client.ID)
		cm.mu.Unlock()
		rm.RemoveClient(client.ID)
		close(client.done) // signals HandleConnectionWrite's range loop to stop, safe: only this goroutine closes it
	}()

	// store the client connection in the connection manager}
	// let the client know, what id it has been assigned by the server
	response := fmt.Sprintf("CLIENT_ID %d", client.ID)

	fmt.Printf("CLIENT_ID %d\n", client.ID)

	// conn.Write([]byte(response))
	client.outChan <- protocol.ServerMsg{Msg: response} // sending message to write to the client's out channel

	// reading the messages sent by client
	scanner := bufio.NewScanner(client.conn) // scan the input from the connection req made

	for scanner.Scan() {
		msg := scanner.Text() // removes the \n, so remember to add a \n at the last so that client can read
		cmd := parseCommand(msg)
		if cmd == nil {
			response = "ERROR invalid command"
			client.outChan <- protocol.ServerMsg{Msg: response} //  write error back to client
			continue
		}

		err := cmd.execute(cm, rm, client.ID)
		if err != nil {
			response = fmt.Sprintf("ERROR %s", err.Error())
			client.outChan <- protocol.ServerMsg{Msg: response} //  write error back to client
		}

		// echo
		// _ = fmt.Sprintf("ACK server responded %s\n", strings.ToUpper(msg))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from client %d : %v\n", client.ID, err)
		client.outChan <- protocol.DisconnectedMsg{Err: err}
	}
}

// responding to the client
func (cm *ConnectionManager) HandleConnectionWrite(client Client) {
	for {
		select {
		case msg := <-client.outChan:
			switch msg := msg.(type) {

			case protocol.ServerMsg:
				_, err := client.conn.Write([]byte(msg.Msg + "\n"))
				if err != nil {
					fmt.Println("write error:", err)
				}

			case protocol.DisconnectedMsg:
				fmt.Println("client disconnected:", msg.Err)
				return
			}
		case <-client.done:
			return
		}
	}
}

func (cm *ConnectionManager) getNextId() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.nextId++
	return cm.nextId
}

func (cm *ConnectionManager) Send(receiverID int, msg string) error {
	cm.mu.RLock()
	client, ok := cm.clients[receiverID]
	cm.mu.RUnlock()
	if !ok {
		return fmt.Errorf("receiver %d not found\n", receiverID)
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	select {
	case client.outChan <- protocol.ServerMsg{Msg: msg}:
		return nil
	case <-client.done:
		return fmt.Errorf("client %d disconnected\n", receiverID)
	}
}

func (cm *ConnectionManager) Broadcast(msg string) error {
	cm.mu.RLock()
	clients := make([]Client, 0, len(cm.clients))
	for _, c := range cm.clients {
		clients = append(clients, c)
	}
	cm.mu.RUnlock()

	for _, client := range clients {
		select {
		case client.outChan <- protocol.ServerMsg{Msg: msg}:
		case <-client.done:
			log.Printf("skipping disconnected client %d during broadcast\n", client.ID)
		}
	}
	return nil
}
