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
	Time           uint32 // dtype??
	LevelTime      uint32 // dtype??
	ServerName     [32]byte
	MapName        [8]byte
	MapTitle       [33]byte
	MapMD5         [16]byte
	ActNum         uint8
	IsZone         uint8
	HttpSource     [MAX_MIRROR_LENGTH]byte
	FileNeeded     [MAXFILENEEDED]byte
}

type PlayerInfoPak struct {
	PacketHeader

	Players [MASTERSERVER_MAXPLAYERS]PlayerInfo
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
