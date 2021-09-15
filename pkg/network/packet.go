package network

const MAXAPPLICATION = 16

type PacketHeader struct {
	Checksum  uint32
	Ack       uint8
	AckReturn uint8

	PacketType Opcode
	Reserved   uint8 // Padding
}

const MAXSERVERNAME = 32
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
	Time           uint64 // dtype??
	LevelTime      uint64 // dtype??
	ServerName     [MAXSERVERNAME]byte
	MapName        [8]byte
	MapTitle       [33]byte
	MapMD5         [16]byte
	ActNum         uint8
	IsZone         uint8
	HttpSource     [MAX_MIRROR_LENGTH]byte
	FileNeeded     [MAXFILENEEDED]byte
}
