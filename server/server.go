package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	Cmd *Command
}

var (
	set = NewClientSet()
)

func (s *Server) Start() {
	defer close(set.CloseChan)
	defer close(set.MessageChan)
	go receiver()
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.Cmd.Server, s.Cmd.Port), nil))
}

func receiver() {
	for {
		select {
		case message := <-set.MessageChan:
			conns := set.Conns(message.User)
			for _, c := range conns {
				err := c.WriteMessage(websocket.TextMessage, []byte(string(message.RawContent)+"update"))
				if err != nil {
					log.Println("write:", err)
				}
			}
		case client := <-set.CloseChan:
			set.Remove(client)
		}
	}
}

//noinspection ALL
func handle(writer http.ResponseWriter, request *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	c, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Panicln("upgrade:", err)
		return
	}
	set.Add(request.URL.Query().Get("user"), c)

	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			set.CloseChan <- &Client{User: request.URL.Query().Get("user"), Conn: c}
			break
		}
		set.MessageChan <- &Message{RawContent: message, User: request.URL.Query().Get("user")}
	}
}
