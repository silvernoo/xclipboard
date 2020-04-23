package server

type Command struct {
	Server, Bind, Port, Key, User string
}

func (a Command) IsServerMode() bool {
	return len(a.Bind) > 0
}