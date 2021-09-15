package main

import "github.com/karashiiro/gokart/pkg/game"

func main() {
	m := game.New(15, "gokart server")
	m.Run()
}
