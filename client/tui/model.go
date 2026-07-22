package tui

import (
	"net"

	"github.com/Shiwang0-0/multiplayertetris/game"
)

type Screen int

type model struct {
	game *game.Game
	conn net.Conn
	myID int

	screen Screen
	width  int
	height int

	// the game state for even the opponent is maintained for each client diffrently i.e no game state is shared and only the moves are shared, and each client replicates that move on their end
	opponents map[int]*game.Game // each opponent id ---> each opponent game
	roomID    string
	isCreator bool // is room creator

	lastError string
	connected bool

	voteSecondsLeft   int // sets by the server message parsing : deadline field
	myVote            string
	activePlayerID    int
	eliminated        bool
	eliminationNotice string

	matchOver bool
	winnerID  int
}

func NewModel(g *game.Game, conn net.Conn) *model {
	return &model{
		game:      g,
		conn:      conn,
		opponents: make(map[int]*game.Game),
		screen:    HomeScreen,
		connected: true,
	}
}
