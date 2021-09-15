package network

type BroadcastConnection struct {
	Connections []Connection
}

var _ Connection = BroadcastConnection{}

func (b BroadcastConnection) Send(data []byte) error {
	for _, conn := range b.Connections {
		if conn == nil {
			continue
		}

		// Ignore error, we should log this
		_ = conn.Send(data)
	}

	return nil
}

func (b *BroadcastConnection) Set(conn Connection) {
	for i, c := range b.Connections {
		if c == nil {
			b.Connections[i] = conn
		}
	}
}

func (b *BroadcastConnection) Unset(conn Connection) {
	for i, c := range b.Connections {
		if c == conn {
			b.Connections[i] = nil
		}
	}
}
