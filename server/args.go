package server

type Command struct {
	Server, Bind, Port, User string
	Key                      []byte
}

func (a Command) IsServerMode() bool {
	return len(a.Bind) > 0
}
