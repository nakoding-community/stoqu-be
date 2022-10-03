package ws

import "github.com/gorilla/websocket"

type Channel struct {
	User       string
	Ws         *websocket.Conn
	Done       chan struct{}
	MsgReceive chan Msg
	MsgSend    chan Msg
}

type Msg struct {
	MsgType int
	Msg     []byte
}

func NewChannel(ws *websocket.Conn) *Channel {
	ch := &Channel{
		Ws:         ws,
		Done:       make(chan struct{}),
		MsgReceive: make(chan Msg),
		MsgSend:    make(chan Msg),
	}
	return ch
}
