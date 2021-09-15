package game

import (
	"github.com/karashiiro/gokart/pkg/network"
)

type player struct {
	name string
	conn network.Connection
}
