package game

import "github.com/karashiiro/gokart/pkg/network"

type Manager struct {
	rooms      []*room
	players    []player
	numPlayers int
	maxPlayers int
	motd       string
	broadcast  network.BroadcastConnection
}

func New(maxPlayers int, motd string) *Manager {
	return &Manager{
		rooms:      nil,
		players:    make([]player, maxPlayers),
		numPlayers: 0,
		maxPlayers: maxPlayers,
		motd:       motd,
		broadcast:  network.BroadcastConnection{},
	}
}

func (m *Manager) Run() {
}

func (m *Manager) tryAddPlayer(p *player) bool {
	// Check if capacity has been reached
	if m.numPlayers >= m.maxPlayers {
		return false
	}

	// Try to add the player to an existing room
	for _, r := range m.rooms {
		if r.tryAddPlayer(p) {
			m.broadcast.Set(p.conn, m.numPlayers)
			m.numPlayers++
			return true
		}
	}

	// Create a new room
	newRoom := &room{
		players:      make([]*player, ROOMMAXPLAYERS),
		numPlayers:   0,
		playerInGame: make([]bool, ROOMMAXPLAYERS),
		state:        "setup",
		broadcast:    network.BroadcastConnection{Connections: make([]network.Connection, ROOMMAXPLAYERS)},
	}
	m.rooms = append(m.rooms, newRoom)
	newRoom.tryAddPlayer(p)

	m.broadcast.Set(p.conn, m.numPlayers)
	m.numPlayers++

	return true
}
