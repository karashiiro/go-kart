package network

type Connection interface {
	Send(data interface{}) error
}
