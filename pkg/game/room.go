package game

type room struct {
	players      []*player
	numPlayers   int
	playerInGame []bool
	state        string
}

func (r *room) tryAddPlayer(p *player) bool {
	// Room is in a match
	if r.state == "playing" {
		return false
	}

	// Room is full, or name is taken
	if r.numPlayers >= 15 || !p.isNameGood(r) {
		return false
	}

	// Add player to room
	r.players[r.numPlayers] = p
	r.playerInGame[r.numPlayers] = false
	r.numPlayers++

	return true
}
