package gamenet

import "github.com/karashiiro/gokart/pkg/doom"

const PACKETVERSION = 0

type PacketHeader struct {
	Checksum  uint32
	Ack       uint8 // If not zero the node asks for acknowledgement, the receiver must resend the ack
	AckReturn uint8 // The return of the ack number

	PacketType Opcode
	Reserved   uint8 // Padding
}

type TicCmd struct {
	ForwardMove int8  // -MAXPLMOVE to MAXPLMOVE (50)
	SideMove    int8  // -MAXPLMOVE to MAXPLMOVE (50)
	AngleTurn   int16 // <<16 for angle delta - saved as 1 byte into demos
	Aiming      int16 // vertical aiming, see G_BuildTicCmd
	Buttons     uint16
	DriftTurn   int16 // SRB2Kart: Used for getting drift turn speed
	Latency     uint8 // Netgames: how many tics ago was this ticcmd generated from this player's end?
}

// Client to server packet
type ClientCmdPak struct {
	PacketHeader

	ClientTic   uint8
	ResendFrom  uint8
	Consistency int16
	Cmd         TicCmd
}

// Splitscreen packet
// WARNING: must have the same format of clientcmd_pak, for more easy use
type Client2CmdPak struct {
	PacketHeader

	ClientTic   uint8
	ResendFrom  uint8
	Consistency int16
	Cmd         TicCmd
	Cmd2        TicCmd
}

// 3P Splitscreen packet
// WARNING: must have the same format of clientcmd_pak, for more easy use
type Client3CmdPak struct {
	PacketHeader

	ClientTic   uint8
	ResendFrom  uint8
	Consistency int16
	Cmd         TicCmd
	Cmd2        TicCmd
	Cmd3        TicCmd
}

// 4P Splitscreen packet
// WARNING: must have the same format of clientcmd_pak, for more easy use
type Client4CmdPak struct {
	PacketHeader

	ClientTic   uint8
	ResendFrom  uint8
	Consistency int16
	Cmd         TicCmd
	Cmd2        TicCmd
	Cmd3        TicCmd
	Cmd4        TicCmd
}

// Server to client packet
// this packet is too large
type ServerTicsPak struct {
	PacketHeader

	StartTic uint8
	NumTics  uint8
	NumSlots uint8      // "Slots filled": Highest player number in use plus one.
	Cmds     [45]TicCmd // Normally [BACKUPTIC][MAXPLAYERS] but too large
}

// Sent to client when all consistency data
// for players has been restored
type ResynchEndPak struct {
	PacketHeader

	RandomSeed uint32

	FlagPlayer [2]int8
	FlagLoose  [2]int32
	FlagFlags  [2]int32
	FlagX      [2]doom.Fixed
	FlagY      [2]doom.Fixed
	FlagZ      [2]doom.Fixed

	InGame  uint32                 // Spectator bit for each player
	CTFTeam [doom.MAXPLAYERS]int32 // Which team? (can't be 1 bit, since in regular Match there are no teams)

	// Resynch game scores and the like all at once
	Score     [doom.MAXPLAYERS]uint32 // Everyone's score
	MareScore [doom.MAXPLAYERS]uint32 // SRB2kart: Battle score
	RealTime  [doom.MAXPLAYERS]doom.Tic
	Laps      [doom.MAXPLAYERS]uint8
}

type ResynchPak struct {
	PacketHeader

	PlayerNum uint8

	PlayerState uint8
	PFlags      uint32
	PAnim       uint8

	Aiming        doom.Angle
	CurrentWeapon uint32
	RingWeapons   uint32

	Powers [doom.NUMPOWERS]uint16

	KartStuff  [doom.NUMKARTSTUFF]int32
	FrameAngle doom.Angle

	Health    int32
	Lives     int8
	Continues int8
	ScoreAdd  uint8
	XtraLife  int8
	Pity      int8

	SkinColor uint8
	Skin      int32

	KartSpeed  uint8
	KartWeight uint8

	CharFlags uint32

	Speed      doom.Fixed
	Jumping    uint8
	SecondJump uint8
	Fly1       uint8
	GlideTime  doom.Tic
	Climbing   uint8
	DeadTimer  int32
	Exiting    doom.Tic
	Homing     uint8
	SkidTime   doom.Tic
	CMomentumX doom.Fixed
	CMomentumY doom.Fixed
	RMomentumX doom.Fixed
	RMomentumY doom.Fixed

	WeaponDelay int32
	TossDelay   int32

	StarPostX     int16
	StarPostY     int16
	StarPostZ     int16
	StarPostNum   int32
	StarPostTime  doom.Tic
	StarPostAngle doom.Angle

	MaxLink         int32
	DashSpeed       doom.Fixed
	DashTime        int32
	AnglePos        doom.Angle
	OldAnglePos     doom.Angle
	BumperTime      doom.Tic
	FlyAngle        int32
	DrillTimer      doom.Tic
	LinkCount       int32
	LinkTimer       doom.Tic
	AnotherFlyAngle int32
	NightsTime      doom.Tic
	DrillMeter      int32
	DrillDelay      uint8
	BonusTime       uint8
	Mare            uint8
	LastSideHit     uint16
	LastLineHit     uint16

	LossTime   doom.Tic
	TimesHit   uint8
	OnConveyor int32

	JoinTime doom.Tic

	SplitscreenIndex uint8

	HasMo uint8

	Angle      doom.Angle
	X          doom.Fixed
	Y          doom.Fixed
	Z          doom.Fixed
	MomentumX  doom.Fixed
	MomentumY  doom.Fixed
	MomentumZ  doom.Fixed
	Friction   doom.Fixed
	MoveFactor doom.Fixed

	Tics     int32
	StateNum doom.StateNum
	Flags    uint32
	Flags2   uint32
	EFlags   uint16

	Radius     doom.Fixed
	Height     doom.Fixed
	Scale      doom.Fixed
	DestScale  doom.Fixed
	ScaleSpeed doom.Fixed
}

type ResynchGotPak struct {
	PacketHeader

	ResynchGot uint8
}

const MAXTEXTCMD = 256

type TextCmdPak struct {
	PacketHeader

	TextCmd [MAXTEXTCMD + 1]byte
}

const MAXFILENEEDED = 915
const MAX_MIRROR_LENGTH = 256

type ServerInfoPak struct {
	PacketHeader

	X255           uint8
	PacketVersion  uint8
	Application    [doom.MAXAPPLICATION]byte
	Version        uint8
	Subversion     uint8
	NumberOfPlayer uint8
	MaxPlayer      uint8
	GameType       uint8
	ModifiedGame   uint8
	CheatsEnabled  uint8
	KartVars       uint8
	FileNeededNum  uint8
	Time           doom.Tic
	LevelTime      doom.Tic
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

	GameTic    doom.Tic
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

type FileTxPak struct {
	PacketHeader

	FileId   uint8
	Position uint32
	Size     uint16
	Data     [0]uint8 // Size is variable using hardware_MAXPACKETLENGTH
}

type ClientConfigPak struct {
	PacketHeader

	X255          uint8
	PacketVersion uint8
	Application   [doom.MAXAPPLICATION]byte
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

type AskInfoPak struct {
	PacketHeader

	Version uint8
	Time    doom.Tic // used for ping evaluation
}

type MSAskInfoPak struct {
	PacketHeader

	ClientAddr [22]byte
	Time       doom.Tic // used for ping evaluation
}

type PlayerConfig struct {
	Name    [doom.MAXPLAYERNAME + 1]byte
	Skin    uint8
	Color   uint8
	PFlags  uint32
	Score   uint32
	CTFTeam uint8
}

type FilesNeededConfigPak struct {
	PacketHeader

	First int32
	Num   uint8
	More  uint8
	Files [MAXFILENEEDED]uint8
}

type FilesNeededNumPak struct {
	PacketHeader

	FilesNeeded int32
}

type PingTablePak struct {
	PacketHeader

	PingTable [doom.MAXPLAYERS + 1]uint32
}
