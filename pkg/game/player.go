package game

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/karashiiro/gokart/pkg/doom"
	"github.com/karashiiro/gokart/pkg/network"
)

type player struct {
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
				p.name = randStringBytes(10)
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

const LETTERBYTES = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = LETTERBYTES[rand.Intn(len(LETTERBYTES))]
	}
	return string(b)
}
