package game

import "github.com/karashiiro/gokart/pkg/network"

const ROOMMAXPLAYERS = 15

type room struct {
	players      []*player
	numPlayers   uint8
	playerInGame []bool
	state        string
	broadcast    *network.BroadcastConnection
}

func (r *room) removePlayer(p *player) {
	playerIdx := -1
	for i := 0; i < len(r.players); i++ {
		if r.players[i] == p {
			playerIdx = i
		}
	}

	if playerIdx == -1 {
		return
	}

	r.broadcast.Unset(p.conn)
	r.playerInGame[playerIdx] = false
	r.players[playerIdx] = r.players[r.numPlayers-1]
	r.players[r.numPlayers-1] = nil

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
	r.players[r.numPlayers] = p
	r.playerInGame[r.numPlayers] = false
	r.broadcast.Set(p.conn)
	r.numPlayers++
	p.room = r

	return true
}
