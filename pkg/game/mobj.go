package game

import "github.com/karashiiro/gokart/pkg/doom"

// MObj represents model data. The server holds this in order to resynch data
// for dead players.
type MObj struct {
	// Info for drawing: position.
	X doom.Fixed
	Y doom.Fixed
	Z doom.Fixed

	Angle doom.Angle // orientation

	MomentumX doom.Fixed
	MomentumY doom.Fixed
	MomentumZ doom.Fixed

	Health int32 // for player this is rings + 1

	Friction   doom.Fixed
	MoveFactor doom.Fixed

	Tics   int32 // state tic counter
	State  doom.StateNum
	Flags  uint32 // flags from mobjinfo tables
	Flags2 uint32 // MF2_ flags
	EFlags uint16 // extra flags

	Radius doom.Fixed
	Height doom.Fixed

	Scale      doom.Fixed
	DestScale  doom.Fixed
	ScaleSpeed doom.Fixed
}
