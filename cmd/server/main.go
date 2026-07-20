package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Shiwang0-0/multiplayertetris/server/networking"
)

var PORT = ":8080"

func main() {

	// listens to the port 8080 for client connection req
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Error listening: ", err)
	}
	defer listener.Close()

	fmt.Println("Server listening to Port", PORT)

	cm := networking.NewConnectionManager()
	rm := networking.NewRoomManager()

	// accept client on a forever loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection req", err)
			continue
		}

		client := cm.Accept(conn) // sure that client exist

		// start reader and writer go routines
		go cm.HandleConnectionRead(client, rm)

		go cm.HandleConnectionWrite(client)
	}
}
