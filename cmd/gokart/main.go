package main

import (
	"log"

	"github.com/karashiiro/gokart/pkg/colortext"
	"github.com/karashiiro/gokart/pkg/game"
)

func main() {
	m, err := game.New(&game.ManagerOptions{
		Port:          5029,
		MaxPlayers:    15,
		Motd:          "gokart server active",
		ServerContext: "",
		ServerName: colortext.
			New().
			AppendTextColored("gokart server", colortext.Cyan).
			Build(),
		KartSpeed: game.KartSpeedNormal,
		GameType:  game.GameTypeRace,
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer m.Close()
	m.Run()
}
