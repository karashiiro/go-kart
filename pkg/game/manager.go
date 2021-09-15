package game

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/karashiiro/gokart/pkg/doom"
	"github.com/karashiiro/gokart/pkg/gamenet"
	"github.com/karashiiro/gokart/pkg/network"
)

type Manager struct {
	port          int
	rooms         []*room
	players       map[network.Connection]*player
	numPlayers    uint8
	maxPlayers    uint8
	motd          string
	serverContext string
	serverName    string
	kartSpeed     KartSpeed
	gameType      GameType
	broadcast     network.BroadcastConnection
	server        *net.UDPConn
}

type ManagerOptions struct {
	Port          int
	MaxPlayers    uint8
	Motd          string
	ServerContext string
	ServerName    string
	KartSpeed     KartSpeed
	GameType      GameType
}

func New(opts *ManagerOptions) (*Manager, error) {
	if len([]byte(opts.Motd)) > 254 {
		return nil, errors.New("motd must be at most 254 bytes")
	}

	if len([]byte(opts.ServerContext)) > 8 {
		return nil, errors.New("server context must be at most 8 bytes")
	}

	return &Manager{
		port:          opts.Port,
		rooms:         nil,
		numPlayers:    0,
		maxPlayers:    opts.MaxPlayers,
		motd:          opts.Motd,
		serverContext: opts.ServerContext,
		serverName:    opts.ServerName,
		kartSpeed:     opts.KartSpeed,
		gameType:      opts.GameType,
		broadcast:     network.BroadcastConnection{Connections: make([]network.Connection, opts.MaxPlayers)},
	}, nil
}

func (m *Manager) Close() {
	m.server.Close()
}

func (m *Manager) Run() {
	server, err := net.ListenUDP("udp", &net.UDPAddr{Port: m.port})
	if err != nil {
		log.Fatalln(err)
	}
	m.server = server

	for {
		data := make([]byte, 1024)
		n, addr, err := m.server.ReadFrom(data)
		if err != nil {
			log.Fatalln(err)
		}
		go m.handleConnection(n, network.NewUDPConnection(m.server, addr), data)
	}
}

func (m *Manager) handleConnection(n int, conn network.Connection, data []byte) {
	header := gamenet.PacketHeader{}
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &header)

	switch header.PacketType {
	case gamenet.PT_ASKINFO:
		askInfo := gamenet.AskInfoPak{}
		buf = bytes.NewReader(data)
		binary.Read(buf, binary.LittleEndian, &askInfo)
		m.sendServerInfo(conn, askInfo.Time)
		m.sendPlayerInfo(conn)
	case gamenet.PT_NODETIMEOUT:
		fallthrough
	case gamenet.PT_CLIENTQUIT:
		m.removePlayer(conn)
	}
}

const SV_SPEEDMASK uint8 = 0x03
const SV_DEDICATED uint8 = 0x40

func (m *Manager) sendServerInfo(conn network.Connection, serverTime uint32) {
	serverInfo := gamenet.ServerInfoPak{
		PacketHeader: gamenet.PacketHeader{
			PacketType: gamenet.PT_SERVERINFO,
		},
		X255:           255,
		PacketVersion:  gamenet.PACKETVERSION,
		Version:        doom.VERSION,
		Subversion:     doom.SUBVERSION,
		Time:           serverTime,
		NumberOfPlayer: m.numPlayers,
		MaxPlayer:      m.maxPlayers,
		GameType:       uint8(m.gameType),
		ModifiedGame:   0,
		CheatsEnabled:  0,
		KartVars:       (uint8(m.kartSpeed) & SV_SPEEDMASK) | SV_DEDICATED,
	}

	copy(serverInfo.ServerName[:], m.serverName)
	copy(serverInfo.Application[:], []byte(doom.SRB2APPLICATION))

	copy(serverInfo.MapName[:], []byte("gokartma"))

	copy(serverInfo.MapTitle[:], []byte("Unknown"))
	serverInfo.MapTitle[32] = 0 // Null out the last byte in case of an overrun

	err := gamenet.SendPacket(conn, &serverInfo)
	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) sendPlayerInfo(conn network.Connection) {
	playerInfo := gamenet.PlayerInfoPak{
		PacketHeader: gamenet.PacketHeader{
			PacketType: gamenet.PT_PLAYERINFO,
		},
	}

	// Send as many players as we can, e.g. min(numPlayers, doom.MASTERSERVER_MAXPLAYERS)
	i := 0
	for _, player := range m.players {
		if i >= int(m.numPlayers) {
			playerInfo.Players[i].Node = 255
			continue
		}

		playerInfo.Players[i].Node = uint8(i)

		copy(playerInfo.Players[i].Name[:], []byte(player.name))
		playerInfo.Players[i].Name[doom.MAXPLAYERNAME] = 0 // Null out the last byte in case of an overrun
	}

	err := gamenet.SendPacket(conn, &playerInfo)
	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) removePlayer(conn network.Connection) {
	var p *player
	var ok bool
	if p, ok = m.players[conn]; !ok {
		return
	}

	if p.room != nil {
		p.room.removePlayer(p)
	}

	m.broadcast.Unset(p)
	delete(m.players, p)

	m.numPlayers--
}

func (m *Manager) tryAddPlayer(p *player) bool {
	// Check if capacity has been reached
	if m.numPlayers >= m.maxPlayers {
		return false
	}

	// Try to add the player to an existing room
	for _, r := range m.rooms {
		if r.tryAddPlayer(p) {
			m.broadcast.Set(p, int(m.numPlayers))
			m.numPlayers++
			return true
		}
	}

	// Create a new room
	newRoom := &room{
		players:      make(map[network.Connection]*player, ROOMMAXPLAYERS),
		numPlayers:   0,
		playerInGame: make(map[network.Connection]bool, ROOMMAXPLAYERS),
		state:        "setup",
		broadcast:    network.BroadcastConnection{Connections: make([]network.Connection, ROOMMAXPLAYERS)},
	}
	m.rooms = append(m.rooms, newRoom)
	newRoom.tryAddPlayer(p)

	m.broadcast.Set(p.conn, int(m.numPlayers))
	m.numPlayers++

	return true
}
