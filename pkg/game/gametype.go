package game

type GameType uint8

const (
	GameTypeRace GameType = iota
	GameTypeBattle
	GameTypeRaceVanilla
	GameTypeBattleVanilla
)
