package networking

import (
	"bufio"
	"net"

	tea "charm.land/bubbletea/v2"
	"github.com/Shiwang0-0/multiplayertetris/client/tui"
	"github.com/Shiwang0-0/multiplayertetris/protocol"
)

func ListenToServer(conn net.Conn, p *tea.Program) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			p.Send(protocol.DisconnectedMsg{Err: err})
			return
		}
		p.Send(tui.ParseResponse(line))
	}
}
