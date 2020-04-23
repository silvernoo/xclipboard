package client

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	"xclipboard/server"
)

type Client struct {
	Cmd  *server.Command
}

var (
	lastText string
	done     = make(chan struct{})
)

func (c *Client) Start() {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", c.Cmd.Server, c.Cmd.Port), RawQuery: fmt.Sprintf("user=%s", c.Cmd.User)}
	log.Printf("connecting to %s", u.String())
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer conn.Close()
	go work(conn)
	receiver(conn)
}

func receiver(conn *websocket.Conn) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			all, err := clipboard.ReadAll()
			if all != lastText {
				if err != nil {
					log.Println(err)
				}
				err = conn.WriteMessage(websocket.TextMessage, []byte(all))
				lastText = all
				if err != nil {
					log.Println(err)
				}
			}
		case <-interrupt:
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func work(conn *websocket.Conn) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		s := string(message)
		if s != lastText {
			err := clipboard.WriteAll(s)
			if err != nil {
				log.Println(err)
			}
			lastText = s
		}
	}
}
