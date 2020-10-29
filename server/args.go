package server

type Command struct {
	Server, Bind, Port, Group string
	Key                       []byte
}

func (a Command) IsServerMode() bool {
	return len(a.Bind) > 0
}
