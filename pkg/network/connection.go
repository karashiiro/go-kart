package network

type Connection interface {
	Send(data []byte) error
}
