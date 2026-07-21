package networking

import (
	"fmt"
	"sync"
)

type RoomManager struct {
	clientToRoom  map[int]int
	roomToClients map[int][]int
	roomCreator   map[int]int
	mu            sync.RWMutex
	matches       map[int]*Match
	matchMu       sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		clientToRoom:  make(map[int]int),
		roomToClients: make(map[int][]int),
		roomCreator:   make(map[int]int),
	}
}

func (rm *RoomManager) GetMatch(roomID int) (*Match, bool) {
	rm.matchMu.RLock()
	defer rm.matchMu.RUnlock()
	m, ok := rm.matches[roomID]
	return m, ok
}

// after awiting for otheres to join
func (rm *RoomManager) StartMatchIfReady(roomID int, cm *ConnectionManager) {
	rm.mu.RLock()
	players := append([]int(nil), rm.roomToClients[roomID]...)
	rm.mu.RUnlock()
	if len(players) < 2 {
		return
	}

	rm.matchMu.Lock()
	if _, exists := rm.matches[roomID]; exists {
		rm.matchMu.Unlock()
		return
	}
	match := NewMatch(roomID, players)
	if rm.matches == nil {
		rm.matches = map[int]*Match{}
	}
	rm.matches[roomID] = match
	rm.matchMu.Unlock()

	match.StartVoting(cm, rm)
}

func (rm *RoomManager) RoomBroadcast(roomID int, message string, cm *ConnectionManager) error {
	rm.mu.RLock()
	clientIDs := append([]int(nil), rm.roomToClients[roomID]...)
	rm.mu.RUnlock()

	// no locking while sending
	for _, clientID := range clientIDs {
		if err := cm.Send(clientID, message); err != nil {
			fmt.Println("Error sending message to client:", clientID)
		}
	}

	return nil
}

func (rm *RoomManager) RoomBroadcastExcept(roomID int, senderID int, message string, cm *ConnectionManager) error {
	rm.mu.RLock()
	clientIDs := append([]int(nil), rm.roomToClients[roomID]...)
	rm.mu.RUnlock()

	// no locking while sending
	for _, clientID := range clientIDs {
		if clientID == senderID {
			continue
		}

		if err := cm.Send(clientID, message); err != nil {
			fmt.Println("Error sending message to client:", clientID)
		}
	}

	return nil
}

func (rm *RoomManager) Join(roomID, senderID int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if _, ok := rm.clientToRoom[senderID]; ok {
		return
	}
	rm.clientToRoom[senderID] = roomID
	rm.roomToClients[roomID] = append(rm.roomToClients[roomID], senderID)
}

func (rm *RoomManager) Leave(roomID, senderID int) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if !rm.clientInRoom(roomID, senderID) { // call the private function, to avoid deadlock of RW mutexes
		return fmt.Errorf("client %d is not a member of room %d\n", senderID, roomID)
	}

	delete(rm.clientToRoom, senderID)

	rm.roomToClients[roomID] = remove(rm.roomToClients[roomID], senderID)

	if len(rm.roomToClients[roomID]) == 0 {
		delete(rm.roomToClients, roomID)
	}

	return nil
}

func (rm *RoomManager) GetRoom(clientID int) (int, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	roomID, ok := rm.clientToRoom[clientID]
	return roomID, ok
}

func (rm *RoomManager) ClientInRoom(roomID, target int) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	return rm.clientInRoom(roomID, target)
}

func (rm *RoomManager) clientInRoom(roomID, target int) bool {
	for _, clientID := range rm.roomToClients[roomID] {
		if clientID == target {
			return true
		}
	}

	return false
}

func (rm *RoomManager) RoomExists(roomID int) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	_, ok := rm.roomToClients[roomID]
	return ok
}

func (rm *RoomManager) RemoveClient(clientID int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	roomID := rm.clientToRoom[clientID]
	rm.roomToClients[roomID] = remove(rm.roomToClients[roomID], clientID)
	if len(rm.roomToClients[roomID]) == 0 {
		delete(rm.roomToClients, roomID)
	}
	delete(rm.clientToRoom, clientID)
}

func remove[T comparable](items []T, target T) []T {
	for i, item := range items {
		if item == target {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
}

func (rm *RoomManager) IsCreator(roomID, clientID int) bool {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return rm.roomCreator[roomID] == clientID
}

func (rm *RoomManager) RoomSize(roomID int) int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	return len(rm.roomToClients[roomID])
}
