package game

import "github.com/karashiiro/gokart/pkg/network"

const ROOMMAXPLAYERS = 15

type room struct {
	players      map[network.Connection]*player
	numPlayers   uint8
	playerInGame map[network.Connection]bool
	state        string
	broadcast    network.BroadcastConnection
}

func (r *room) removePlayer(p *player) {
	var ok bool
	if _, ok = r.players[p]; !ok {
		return
	}

	r.broadcast.Unset(p)
	delete(r.players, p)
	delete(r.playerInGame, p)

	r.numPlayers--
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
	r.players[p] = p
	r.playerInGame[p] = false
	r.broadcast.Set(p, int(r.numPlayers))
	r.numPlayers++
	p.room = r

	return true
}
