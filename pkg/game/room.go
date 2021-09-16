package game

import (
	"log"

	"github.com/karashiiro/gokart/pkg/doom"
	"github.com/karashiiro/gokart/pkg/gamenet"
	"github.com/karashiiro/gokart/pkg/network"
)

const ROOMMAXPLAYERS = 15

type room struct {
	players      []*player
	numPlayers   uint8
	playerInGame []bool
	state        string // This is a made-up placeholder, TODO
	broadcast    *network.BroadcastConnection
}

func (r *room) handlePacketFromPlayer(p *player, data []byte) {
	header := gamenet.PacketHeader{}
	gamenet.ReadPacket(data, &header)

	log.Printf("Got packet from %s with type %d", p.name, header.PacketType)
}

// isTicCmdHacked returns true if speedhacking is detected
func (r *room) isTicCmdHacked(cmd *gamenet.TicCmd) bool {
	if cmd.ForwardMove > doom.MAXPLMOVE || cmd.ForwardMove < -doom.MAXPLMOVE ||
		cmd.SideMove > doom.MAXPLMOVE || cmd.SideMove < -doom.MAXPLMOVE ||
		cmd.DriftTurn > doom.KART_FULLTURN || cmd.DriftTurn < -doom.KART_FULLTURN {
		return true
	}

	return false
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
