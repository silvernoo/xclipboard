package server

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	RawContent []byte
	User       string
}

type Client struct {
	Conn *websocket.Conn
	User string
	Cmd  *Command
}

type ClientSet struct {
	MessageChan chan *Message
	CloseChan   chan *Client
	conns       map[string][]*websocket.Conn
}

func NewClientSet() *ClientSet {
	var set ClientSet
	set.conns = make(map[string][]*websocket.Conn)
	set.CloseChan = make(chan *Client)
	set.MessageChan = make(chan *Message)
	return &set
}

func (s *ClientSet) Add(get string, c *websocket.Conn) {
	s2 := s.conns[get]
	s2 = append(s2, c)
	s.conns[get] = s2
}

func (s ClientSet) Conns(get string) []*websocket.Conn {
	return s.conns[get]
}

func (s *ClientSet) Remove(client *Client) {
	conns := s.conns[client.User]
	var i int
	for i = 0; i < len(conns); i++ {
		if conns[i] == client.Conn {
			break
		}
	}
	s.conns[client.User] = append(conns[:i], conns[i+1:]...)
}
