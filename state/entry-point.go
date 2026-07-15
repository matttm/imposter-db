package state

type EntryPoint struct {
}

func (e EntryPoint) IsComplete() bool {
	return false
}

func (e EntryPoint) Handle(f *uint32, schema string, remote net.Conn, client net.Conn, username, password string, cancel context.CancelFunc) state.State {
	// Implement the handling logic here
	return e
}
