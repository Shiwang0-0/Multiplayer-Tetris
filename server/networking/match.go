package networking

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Phase int

const (
	PhaseWaiting Phase = iota
	PhaseVoting
	PhasePlaying
)

var pieceTypes = []string{"I", "O", "T", "L"}

type Match struct {
	roomID          int
	players         []int // turn order, fixed once the match starts
	turnIdx         int
	phase           Phase
	votes           map[int]string // playerID -> chosen piece
	deadLineSeconds int
	mu              sync.Mutex
}

func NewMatch(roomID int, players []int) *Match {
	return &Match{
		roomID:          roomID,
		players:         players,
		deadLineSeconds: 10,
		votes:           map[int]string{},
	}
}

func (mt *Match) OnLocked(playerID int, cm *ConnectionManager, rm *RoomManager) error {
	mt.mu.Lock()
	if mt.phase != PhasePlaying || playerID != mt.activePlayer() {
		mt.mu.Unlock()
		return fmt.Errorf("not your turn to lock")
	}
	mt.turnIdx = (mt.turnIdx + 1) % len(mt.players)
	mt.mu.Unlock()
	mt.StartVoting(cm, rm) // if not PhasePlaying and not active player, start the voting for the next player
	return nil
}

func (mt *Match) StartVoting(cm *ConnectionManager, rm *RoomManager) {
	mt.mu.Lock()
	mt.phase = PhaseVoting
	mt.votes = map[int]string{}
	active := mt.activePlayer()
	deadline := mt.deadLineSeconds
	mt.mu.Unlock()

	rm.RoomBroadcast(mt.roomID, fmt.Sprintf("VOTING_START %d %d", active, deadline), cm)

	// voting time for 10 sec, server will boradcast the start message once the voting timeout ends
	time.AfterFunc(10*time.Second, func() {
		mt.resolveVotes(cm, rm)
	})
}

func (mt *Match) SubmitVote(voterID int, piece string) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if mt.phase != PhaseVoting {
		return fmt.Errorf("not currently voting")
	}
	if voterID == mt.activePlayer() {
		return fmt.Errorf("the active player can't vote for their own piece")
	}
	valid := false
	for _, p := range pieceTypes {
		if p == piece {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid piece %q", piece)
	}
	mt.votes[voterID] = piece
	return nil
}

func (mt *Match) resolveVotes(cm *ConnectionManager, rm *RoomManager) {
	mt.mu.Lock()
	if mt.phase != PhaseVoting {
		mt.mu.Unlock()
		return // already resolved or match moved on
	}
	freqPerPiece := map[string]int{}
	for _, piece := range mt.votes {
		freqPerPiece[piece]++
	}
	winner, best := "", -1
	for piece, count := range freqPerPiece {
		if count > best {
			best, winner = count, piece
		}
	}
	if winner == "" {
		winner = pieceTypes[rand.Intn(len(pieceTypes))] // nobody voted — pick randomly rather than stall
	}
	active := mt.activePlayer()
	mt.phase = PhasePlaying
	mt.mu.Unlock()

	rm.RoomBroadcast(mt.roomID, fmt.Sprintf("TURN_START %d %s", active, winner), cm)
}

func (m *Match) activePlayer() int {
	return m.players[m.turnIdx]
}
