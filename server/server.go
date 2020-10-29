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
	defer set.Close()
	go receiver()
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", s.Cmd.Bind, s.Cmd.Port), nil))
}

func receiver() {
	for {
		select {
		case message := <-set.MessageChan:
			conns := set.Conns(message.Group)
			fmt.Println(len(conns))
			for _, c := range conns {
				err := c.WriteMessage(websocket.TextMessage, []byte(string(message.RawContent)))
				if err != nil {
					log.Println("write:", err)
				}
			}
		case client := <-set.CloseChan:
			set.Remove(client)
		}
	}
}

func handle(writer http.ResponseWriter, request *http.Request) {
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	c, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Panicln("upgrade:", err)
		return
	}
	set.Add(request.URL.Query().Get("group"), c)
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			set.CloseChan <- &Client{Group: request.URL.Query().Get("group"), Conn: c}
			break
		}
		set.MessageChan <- &Message{RawContent: message, Group: request.URL.Query().Get("group")}
	}
}
