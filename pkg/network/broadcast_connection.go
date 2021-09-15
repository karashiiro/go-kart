package network

type BroadcastConnection struct {
	Connections []Connection
}

func (b BroadcastConnection) Send(data interface{}) error {
	for _, conn := range b.Connections {
		// Ignore error, we should log this
		_ = conn.Send(data)
	}

	return nil
}

func (b *BroadcastConnection) Set(conn Connection, i int) {
	b.Connections[i] = conn
}

var _ Connection = BroadcastConnection{}
