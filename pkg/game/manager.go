package game

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/karashiiro/gokart/pkg/doom"
	"github.com/karashiiro/gokart/pkg/gamenet"
	"github.com/karashiiro/gokart/pkg/network"
	"github.com/karashiiro/gokart/pkg/text"
)

type Manager struct {
	port          int
	rooms         []*room
	players       map[string]*player
	numPlayers    uint8
	maxPlayers    uint8
	motd          string
	serverContext string
	serverName    string
	kartSpeed     KartSpeed
	gameType      GameType
	broadcast     *network.BroadcastConnection
	server        *net.UDPConn
}

type ManagerOptions struct {
	Port       int
	MaxPlayers uint8
	Motd       string
	ServerName string
	KartSpeed  KartSpeed
	GameType   GameType
}

func New(opts *ManagerOptions) (*Manager, error) {
	if len([]byte(opts.Motd)) > 254 {
		return nil, errors.New("motd must be at most 254 bytes")
	}

	return &Manager{
		port:          opts.Port,
		rooms:         nil,
		numPlayers:    0,
		maxPlayers:    opts.MaxPlayers,
		motd:          opts.Motd,
		serverContext: text.RandStringBytes(8),
		serverName:    opts.ServerName,
		kartSpeed:     opts.KartSpeed,
		gameType:      opts.GameType,
		broadcast:     network.NewBroadcastConnection(int(opts.MaxPlayers)),
	}, nil
}

func (m *Manager) Close() {
	shutdown := gamenet.PacketHeader{
		PacketType: gamenet.PT_SERVERSHUTDOWN,
	}
	gamenet.SendPacket(m.broadcast, &shutdown)
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

func (m *Manager) handleConnect(conn network.Connection) {
	if m.numPlayers >= m.maxPlayers {
		m.sendRefuse(conn, fmt.Sprintf("Maximum players reached: %d", m.maxPlayers))
	}
}

func (m *Manager) sendRefuse(conn network.Connection, reason string) {
	refuse := gamenet.ServerRefusePak{
		PacketHeader: gamenet.PacketHeader{
			PacketType: gamenet.PT_SERVERREFUSE,
		},
	}

	copy(refuse.Reason[:], []byte(reason))

	err := gamenet.SendPacket(conn, &refuse)
	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) removePlayer(conn network.Connection) {
	var p *player
	var ok bool
	if p, ok = m.players[conn.Addr().String()]; !ok {
		return
	}

	if p.room != nil {
		p.room.removePlayer(p)
	}

	m.broadcast.Unset(conn)
	delete(m.players, conn.Addr().String())

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
			m.broadcast.Set(p.conn)
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
		broadcast:    network.NewBroadcastConnection(ROOMMAXPLAYERS),
	}
	m.rooms = append(m.rooms, newRoom)
	newRoom.tryAddPlayer(p)

	m.broadcast.Set(p.conn)
	m.numPlayers++

	return true
}
