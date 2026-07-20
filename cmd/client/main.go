package main

import (
	"log"
	"net"

	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/client/networking"
	"github.com/Shiwang0-0/multiplayertetris/client/tui"
	"github.com/Shiwang0-0/multiplayertetris/game"
)

func main() {
	// set up the connection
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal("Error connecting to server: ", err)
	}
	defer conn.Close()

	// initalize a game

	g := game.NewGame()

	tuiModel := tui.NewModel(g, conn) // sends the commands pressed over the conn

	// initalizing the bubble tea program
	p := tea.NewProgram(tuiModel)

	// this dedicated go routines, listens to the messages from the server, and then feeds it to the bubble tea running program
	go networking.ListenToServer(conn, p)

	// now run the bubble tea program
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
