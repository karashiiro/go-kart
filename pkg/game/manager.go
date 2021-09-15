package game

type Manager struct {
	rooms      []room
	players    int
	maxPlayers int
}

func New(maxPlayers int) *Manager {
	return &Manager{maxPlayers: maxPlayers}
}

func (m *Manager) Run() {
}
