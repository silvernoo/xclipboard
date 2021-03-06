package main

import (
	"flag"
	"xclipboard/client"
	"xclipboard/server"
)

// xclipboard -s 127.0.0.1 -p 9000 -k 123 -u home
// xclipboard -b 127.0.0.1 -p 9000
func main() {
	var cmd server.Command
	var key string
	flag.StringVar(&cmd.Server, "s", "", "server address")
	flag.StringVar(&cmd.Bind, "b", "", "binding address")
	flag.StringVar(&cmd.Port, "p", "9000", "binding address")
	flag.StringVar(&key, "k", "&*……UJM·12", "encrypt key")
	flag.StringVar(&cmd.Group, "g", "default", "group")
	flag.Parse()
	bytes := make([]byte, 32)
	copy(bytes, key)
	cmd.Key = bytes
	if cmd.IsServerMode() {
		(&server.Server{Cmd: &cmd}).Start()
	} else {
		(&client.Client{Cmd: &cmd}).Start()
	}
}
