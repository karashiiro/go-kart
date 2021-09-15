package game

import "github.com/karashiiro/gokart/pkg/network"

const ROOMMAXPLAYERS = 15

type room struct {
	players      []*player
	numPlayers   uint8
	playerInGame []bool
	state        string
	broadcast    network.BroadcastConnection
}

func (r *room) tryAddPlayer(p *player) bool {
	// Room is in a match
	if r.state == "playing" {
		return false
	}

	// Room is full, or name is taken
	if r.numPlayers >= ROOMMAXPLAYERS || !p.isNameGood(r) {
		return false
	}

	// Add player to room
	r.players[r.numPlayers] = p
	r.playerInGame[r.numPlayers] = false
	r.broadcast.Set(p.conn, int(r.numPlayers))
	r.numPlayers++

	return true
}
