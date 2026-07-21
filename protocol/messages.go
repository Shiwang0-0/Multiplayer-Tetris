package protocol

// Msg is anything sent server -> client over the wire.
// The unexported isMsg() method seals this interface: only types defined
// in this package can implement Msg, so a stray struct from elsewhere
// can never accidentally satisfy it.
type Msg interface {
	isMsg()
}

type ServerMsg struct {
	Msg string
}

func (ServerMsg) isMsg() {}

type ServerErrorMsg struct {
	Msg string
}

func (ServerErrorMsg) isMsg() {}

type ClientIDMsg struct {
	ID int
}

func (ClientIDMsg) isMsg() {}

type DisconnectedMsg struct {
	Err error
}

func (DisconnectedMsg) isMsg() {}

type JoinedMsg struct {
	RoomID   int
	ClientID int
}

func (JoinedMsg) isMsg() {}

type OpponentMoveMsg struct {
	SenderID int
	Move     string
}

func (OpponentMoveMsg) isMsg() {}

type VotingStartMsg struct {
	ActivePlayerID  int
	DeadlineSeconds int
}

func (VotingStartMsg) isMsg() {}

type TurnStartMsg struct {
	ActivePlayerID int
	Piece          string // "I", "O", "T", "L"
}

func (TurnStartMsg) isMsg() {}
