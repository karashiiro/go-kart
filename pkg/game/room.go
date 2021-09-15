package game

type room struct {
	players      []*player
	playerInGame []bool
	state        string
}
