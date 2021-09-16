package game

import (
	"strconv"
	"strings"

	"github.com/karashiiro/gokart/pkg/doom"
	"github.com/karashiiro/gokart/pkg/network"
	"github.com/karashiiro/gokart/pkg/text"
)

type player struct {
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

	mobj *MObj
	name string
	conn network.Connection
	room *room
}

func (p player) Send(data []byte) error {
	return p.conn.Send(data)
}

func (p *player) isNameGood(r *room) bool {
	// Empty or too long
	if len(p.name) == 0 || len(p.name) > doom.MAXPLAYERNAME {
		return false
	}

	// Starts/ends with a space
	if p.name[0] == ' ' || p.name[len(p.name)-1] == ' ' {
		return false
	}

	// Starts with a digit
	if p.name[0] >= 48 && p.name[0] < 58 {
		return false
	}

	// Starts with an admin symbol
	if p.name[0] == '@' || p.name[0] == '~' {
		return false
	}

	// Check if it contains a non-printing character.
	// Note: ANSI C isprint() considers space a printing character.
	// Also don't allow semicolons, since they are used as
	// console command separators.

	// Also, anything over 0x80 is disallowed too, since compilers love to
	// differ on whether they're printable characters or not.
	for _, c := range p.name {
		if !strconv.IsPrint(c) || c == ';' || c >= 0x80 {
			return false
		}
	}

	// Check if a player is currently using the name, case-insensitively
	for i, otherPlayer := range r.players {
		if otherPlayer != nil && otherPlayer != p && r.playerInGame[i] && strings.EqualFold(otherPlayer.name, p.name) {
			// We shouldn't kick people out just because
			// they joined the game with the same name
			// as someone else -- modify the name instead.

			// Slowly strip characters off the end of the
			// name until we no longer have a duplicate.
			if len(p.name) > 1 {
				p.name = p.name[:len(p.name)-1]
				if !p.isNameGood(r) {
					return false
				}
			} else if len(p.name) == 1 {
				// last ditch effort
				p.name = text.RandStringBytes(10)
				if !p.isNameGood(r) {
					return false
				}
			} else {
				// Nothing worked, kick them :(
				return false
			}
		}
	}

	return true
}
