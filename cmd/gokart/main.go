package main

import (
	"log"

	"github.com/karashiiro/gokart/pkg/game"
)

func main() {
	m, err := game.New(&game.ManagerOptions{
		Port:          5029,
		MaxPlayers:    15,
		Motd:          "gokart server active",
		ServerContext: "",
		ServerName:    string([]byte{0x88}) + "gokart server",
		KartSpeed:     game.KartSpeedNormal,
		GameType:      game.GameTypeRace,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer m.Close()
	m.Run()
}
