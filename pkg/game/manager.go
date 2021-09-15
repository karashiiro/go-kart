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
)

type Manager struct {
	port          int
	rooms         []*room
	players       []player
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
		players:       make([]player, opts.MaxPlayers),
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
		n, c, err := m.server.ReadFrom(data)
		if err != nil {
			log.Fatalln(err)
		}
		go m.handleConnection(n, c, data)
	}
}

func (m *Manager) handleConnection(n int, addr net.Addr, data []byte) {
	header := gamenet.PacketHeader{}
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &header)

	switch header.PacketType {
	case gamenet.PT_ASKINFO:
		askInfo := gamenet.AskInfoPak{}
		buf = bytes.NewReader(data)
		binary.Read(buf, binary.LittleEndian, &askInfo)
		m.sendServerInfo(addr, askInfo.Time)
		m.sendPlayerInfo(addr)
	}
}

const SV_SPEEDMASK uint8 = 0x03
const SV_DEDICATED uint8 = 0x40

func (m *Manager) sendServerInfo(addr net.Addr, serverTime uint32) {
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

	m.sendPacket(addr, &serverInfo)
}

func (m *Manager) sendPlayerInfo(addr net.Addr) {
	playerInfo := gamenet.PlayerInfoPak{
		PacketHeader: gamenet.PacketHeader{
			PacketType: gamenet.PT_PLAYERINFO,
		},
	}

	// Send as many players as we can, e.g. min(numPlayers, doom.MASTERSERVER_MAXPLAYERS)
	for i, player := range playerInfo.Players {
		if i >= int(m.numPlayers) {
			player.Node = 255
			continue
		}

		player.Node = uint8(i)

		copy(player.Name[:], []byte(m.players[i].name))
		player.Name[doom.MAXPLAYERNAME] = 0 // Null out the last byte in case of an overrun
	}

	m.sendPacket(addr, &playerInfo)
}

func (m *Manager) sendPacket(addr net.Addr, data interface{}) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)

	sendBuf := buf.Bytes()
	checksum := netBufferChecksum(sendBuf[4:])
	binary.LittleEndian.PutUint32(sendBuf[0:4], checksum)

	_, err := m.server.WriteTo(sendBuf, addr)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
}

func netBufferChecksum(buf []byte) uint32 {
	checksum := uint32(0x1234567)
	length := len(buf)
	for i := 0; i < length; i++ {
		checksum += uint32(buf[i]) * uint32(i+1)
	}
	return checksum
}

func (m *Manager) tryAddPlayer(p *player) bool {
	// Check if capacity has been reached
	if m.numPlayers >= m.maxPlayers {
		return false
	}

	// Try to add the player to an existing room
	for _, r := range m.rooms {
		if r.tryAddPlayer(p) {
			m.broadcast.Set(p.conn, int(m.numPlayers))
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

	m.broadcast.Set(p.conn, int(m.numPlayers))
	m.numPlayers++

	return true
}
