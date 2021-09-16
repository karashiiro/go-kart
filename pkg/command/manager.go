package command

import (
	"fmt"

	"github.com/karashiiro/gokart/pkg/network"
)

type NetXCmdDelegate func(conn network.Connection, params []interface{}) error

type Manager struct {
	commands map[NetXCmd]NetXCmdDelegate
}

func (m *Manager) RegisterNetXCommand(id NetXCmd, del NetXCmdDelegate) error {
	if id > MAXNETXCMD {
		return fmt.Errorf("command id %d too big", id)
	}

	if _, ok := m.commands[id]; ok {
		return fmt.Errorf("command id %d already used", id)
	}

	m.commands[id] = del

	return nil
}

func (m *Manager) ExecuteNetXCommand(id NetXCmd, conn network.Connection, params []interface{}) error {
	if cmd, ok := m.commands[id]; ok {
		return cmd(conn, params)
	}

	return fmt.Errorf("command id %d not registered", id)
}
