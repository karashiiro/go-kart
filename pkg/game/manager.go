package game

type Manager struct {
	rooms      []room
	players    []player
	numPlayers int
	maxPlayers int
}

func New(maxPlayers int) *Manager {
	return &Manager{
		rooms:      nil,
		players:    make([]player, maxPlayers),
		numPlayers: 0,
		maxPlayers: maxPlayers,
	}
}

func (m *Manager) Run() {
}
