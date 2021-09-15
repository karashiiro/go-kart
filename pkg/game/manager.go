package game

import "github.com/karashiiro/gokart/pkg/network"

type Manager struct {
	rooms      []*room
	players    []player
	numPlayers int
	maxPlayers int
	motd       string
	broadcast  network.Connection
}

func New(maxPlayers int, motd string) *Manager {
	return &Manager{
		rooms:      nil,
		players:    make([]player, maxPlayers),
		numPlayers: 0,
		maxPlayers: maxPlayers,
		motd:       motd,
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
			return true
		}
	}

	// Create a new room
	newRoom := &room{
		players:      make([]*player, 15),
		numPlayers:   0,
		playerInGame: make([]bool, 15),
		state:        "setup",
	}
	newRoom.tryAddPlayer(p)
	m.rooms = append(m.rooms, newRoom)

	return true
}
