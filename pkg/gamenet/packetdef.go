package gamenet

import "github.com/karashiiro/gokart/pkg/doom"

const PACKETVERSION = 0

const MAXAPPLICATION = 16

type PacketHeader struct {
	Checksum  uint32
	Ack       uint8
	AckReturn uint8

	PacketType Opcode
	Reserved   uint8 // Padding
}

type AskInfoPak struct {
	PacketHeader

	Version uint8
	Time    uint32
}

const MAXFILENEEDED = 915
const MAX_MIRROR_LENGTH = 256

type ServerInfoPak struct {
	PacketHeader

	X255           uint8
	PacketVersion  uint8
	Application    [MAXAPPLICATION]byte
	Version        uint8
	Subversion     uint8
	NumberOfPlayer uint8
	MaxPlayer      uint8
	GameType       uint8
	ModifiedGame   uint8
	CheatsEnabled  uint8
	KartVars       uint8
	FileNeededNum  uint8
	Time           uint32
	LevelTime      uint32
	ServerName     [32]byte
	MapName        [8]byte
	MapTitle       [33]byte
	MapMD5         [16]byte
	ActNum         uint8
	IsZone         uint8
	HttpSource     [MAX_MIRROR_LENGTH]byte
	FileNeeded     [MAXFILENEEDED]byte
}

type ServerConfigPak struct {
	PacketHeader

	Version    uint8
	Subversion uint8

	ServerPlayer uint8
	TotalSlotNum uint8 // "Slots": highest player number in use plus one.

	GameTic    uint32
	ClientNode uint8
	GameState  uint8

	// 0xFF == not in game; else player skin num
	PlayerSkins [doom.MAXPLAYERS]uint8
	PlayerColor [doom.MAXPLAYERS]uint8

	GameType     uint8
	ModifiedGame uint8
	AdminPlayers [doom.MAXPLAYERS]int8 // Needs to be signed

	ServerContext [8]byte // Unique context id

	// Discord info (always defined for net compatibility)
	MaxPlayer      uint8
	AllowNewPlayer bool
	DiscordInvites bool

	VarLengthInputs [0]uint8 // Playernames and netvars
}

type ClientConfigPak struct {
	PacketHeader

	X255          uint8
	PacketVersion uint8
	Application   [MAXAPPLICATION]byte
	Version       uint8
	Subversion    uint8
	LocalPlayers  uint8 // number of splitscreen players
	Mode          uint8
}

type PlayerInfoPak struct {
	PacketHeader

	Players [doom.MASTERSERVER_MAXPLAYERS]PlayerInfo
}

// PlayerInfo represents shorter player information for external use.
type PlayerInfo struct {
	Node            uint8
	Name            [doom.MAXPLAYERNAME + 1]byte
	Reserved        [4]uint8
	Team            uint8
	Skin            uint8
	Data            uint8 // Color is first four bits, hasflag, isit and issuper have one bit each, the last is unused.
	Score           uint32
	SecondsInServer uint16
}

type ServerRefusePak struct {
	PacketHeader

	Reason [255]byte
}
