package main

import (
	"flag"
	"xclipboard/client"
	"xclipboard/server"
)

// xclipboard -s 127.0.0.1 -p 9000 -k 123 -u home
// xclipboard -b 127.0.0.1 -p 9000 -k 123
func main() {
	var cmd server.Command
	flag.StringVar(&cmd.Server, "s", "", "server address")
	flag.StringVar(&cmd.Bind, "b", "", "binding address")
	flag.StringVar(&cmd.Port, "p", "9000", "binding address")
	flag.StringVar(&cmd.Key, "k", "2673890", "encrypt key")
	flag.StringVar(&cmd.User, "u", "default", "user")
	flag.Parse()
	if cmd.IsServerMode() {
		s := server.Server{Cmd: &cmd}
		s.Start()
	} else {
		c := client.Client{Cmd: &cmd}
		c.Start()
	}
}
