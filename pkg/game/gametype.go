package game

type GameType uint8

const (
	GameTypeRace GameType = iota + 2
	GameTypeBattle
)
