package main

import (
	"log"

	"github.com/karashiiro/gokart/pkg/game"
)

func main() {
	m, err := game.New(15, "gokart server", "")
	if err != nil {
		log.Fatalln(err)
	}
	m.Run()
}
