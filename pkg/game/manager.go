package game

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/karashiiro/gokart/pkg/command"
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
	commands      *command.Manager
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
		players:       make(map[string]*player),
		numPlayers:    0,
		maxPlayers:    opts.MaxPlayers,
		motd:          opts.Motd,
		serverContext: text.RandStringBytes(8),
		serverName:    opts.ServerName,
		kartSpeed:     opts.KartSpeed,
		gameType:      opts.GameType,
		broadcast:     network.NewBroadcastConnection(int(opts.MaxPlayers)),
		commands:      &command.Manager{},
	}, nil
}

func (m *Manager) Close() {
	if m.server == nil {
		return
	}

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
		// PT_SERVERINFO should be the largest packet at 1024 bytes.
		// d_clisrv.h notes 64kB packets under doomdata_t, but those
		// are probably junk numbers.
		data := make([]byte, 1024)
		_, addr, err := m.server.ReadFrom(data)
		if err != nil {
			log.Fatalln(err)
		}

		var p *player
		var ok bool
		if p, ok = m.players[addr.String()]; ok {
			// Player messages are handled by the room they're in
			if p.room != nil {
				go p.room.handlePacketFromPlayer(p, data)
			}
		} else {
			go m.handlePacketFromAwayNode(network.NewUDPConnection(m.server, addr), data)
		}
	}
}

// handlePacketFromAwayNode handles messages from connections that aren't
// necessarily players. This is used for things such as getting server info
// for directories, and getting addon info.
func (m *Manager) handlePacketFromAwayNode(conn network.Connection, data []byte) {
	header := gamenet.PacketHeader{}
	gamenet.ReadPacket(data, &header)

	log.Printf("Got packet from %s with type %d", conn.Addr().String(), header.PacketType)

	switch header.PacketType {
	case gamenet.PT_ASKINFO:
		askInfo := gamenet.AskInfoPak{}
		gamenet.ReadPacket(data, &askInfo)
		m.sendServerInfo(conn, askInfo.Time)
		m.sendPlayerInfo(conn)
	case gamenet.PT_NODETIMEOUT:
		fallthrough
	case gamenet.PT_CLIENTQUIT:
		m.removePlayer(conn)
	case gamenet.PT_CLIENTJOIN:
		cfg := gamenet.ClientConfigPak{}
		gamenet.ReadPacket(data, &cfg)
		m.handleConnect(conn, &cfg)
	}
}

func (m *Manager) handleConnect(conn network.Connection, cfg *gamenet.ClientConfigPak) {
	if cfg.X255 != 255 || cfg.PacketVersion != gamenet.PACKETVERSION {
		m.sendRefuse(conn, "Incompatible packet formats.")
	} else if !strings.EqualFold(string(cfg.Application[:]), string(doom.SRB2APPLICATION)) {
		m.sendRefuse(conn, "Different SRB2 modifications\nare not compatible.")
	} else if cfg.Version != doom.VERSION || cfg.Subversion != doom.SUBVERSION {
		m.sendRefuse(conn, fmt.Sprintf("Different SRB2Kart versions cannot\nplay a netgame!\n(server version %d.%d)", doom.VERSION, doom.SUBVERSION))
	} else if m.numPlayers+cfg.LocalPlayers > m.maxPlayers {
		m.sendRefuse(conn, fmt.Sprintf("Number of local players\nwould exceed maximum: %d", m.maxPlayers))
	} else if m.numPlayers >= m.maxPlayers {
		m.sendRefuse(conn, fmt.Sprintf("Maximum players reached: %d", m.maxPlayers))
	} else if cfg.LocalPlayers > 4 {
		// We only have packet types for up to 4 splitscreen players
		m.sendRefuse(conn, "Too many players from\nthis node.")
	} else if cfg.LocalPlayers == 0 { // Stealth join?
		m.sendRefuse(conn, "No players from\nthis node.")
	} else {
		// Client authorized to join
		p := player{conn: conn}
		if m.tryAddPlayer(&p) {
			m.players[p.conn.Addr().String()] = &p
		}

		m.sendServerConfig(p.conn)
	}
}

const SV_SPEEDMASK uint8 = 0x03
const SV_DEDICATED uint8 = 0x40

// sendServerInfo sends the server info. This is used to populate
// server directories.
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
		ModifiedGame:   0,
		CheatsEnabled:  0,
		KartVars:       (uint8(m.kartSpeed) & SV_SPEEDMASK) | SV_DEDICATED,
	}

	// SRB2Kart: Vanilla's gametype constants for MS support
	if m.gameType == GameTypeRace {
		serverInfo.GameType = uint8(GameTypeRaceVanilla)
	} else {
		serverInfo.GameType = uint8(GameTypeBattleVanilla)
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

// sendPlayerInfo sends information about the players on this server.
// This is used to populate server directories.
func (m *Manager) sendPlayerInfo(conn network.Connection) {
	playerInfo := gamenet.PlayerInfoPak{
		PacketHeader: gamenet.PacketHeader{
			PacketType: gamenet.PT_PLAYERINFO,
		},
	}

	// Send as many players as we can, e.g. min(numPlayers, doom.MASTERSERVER_MAXPLAYERS)
	i := 0
	for _, player := range m.players {
		if i >= doom.MASTERSERVER_MAXPLAYERS {
			break
		}

		if i >= int(m.numPlayers) {
			playerInfo.Players[i].Node = 255
			continue
		}

		playerInfo.Players[i].Node = uint8(i)

		copy(playerInfo.Players[i].Name[:], []byte(player.name))
		playerInfo.Players[i].Name[doom.MAXPLAYERNAME] = 0 // Null out the last byte in case of an overrun

		i++
	}

	err := gamenet.SendPacket(conn, &playerInfo)
	if err != nil {
		log.Println(err)
	}
}

func (m *Manager) sendServerConfig(conn network.Connection) {
	serverConfig := gamenet.ServerConfigPak{
		PacketHeader: gamenet.PacketHeader{
			PacketType: gamenet.PT_SERVERCFG,
		},
		Version:        doom.VERSION,
		Subversion:     doom.SUBVERSION,
		ServerPlayer:   m.maxPlayers,
		TotalSlotNum:   m.numPlayers,
		GameTic:        uint32(time.Now().Unix()),
		ClientNode:     0,
		GameState:      0,
		GameType:       uint8(m.gameType),
		ModifiedGame:   0,
		MaxPlayer:      m.maxPlayers,
		AllowNewPlayer: true,
		DiscordInvites: false,
	}

	// Send as many players as we can, e.g. min(numPlayers, doom.MAXPLAYERS)
	i := 0
	for _, player := range m.players {
		if i >= doom.MAXPLAYERS {
			break
		}

		// TODO: Implement admin players
		serverConfig.AdminPlayers[i] = -1

		if i >= int(m.numPlayers) {
			serverConfig.PlayerSkins[i] = 0xFF
			serverConfig.PlayerColor[i] = 0xFF
			continue
		}

		serverConfig.PlayerSkins[i] = uint8(player.Skin)
		serverConfig.PlayerColor[i] = player.SkinColor

		i++
	}

	copy(serverConfig.ServerContext[:], m.serverContext)

	err := gamenet.SendPacket(conn, &serverConfig)
	if err != nil {
		log.Println(err)
	}
}

// sendRefuse denies a connection's join request, sending
// back the provided reason.
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

// removePlayer removes a player from the server.
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

// tryAddPlayer attempts to add a player to this room,
// returning true if the operation was successful.
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
