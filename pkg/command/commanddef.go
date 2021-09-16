package command

type NetXCmd = int

const (
	XD_NAMEANDCOLOR NetXCmd = iota + 1
	XD_WEAPONPREF           // 2
	XD_KICK                 // 3
	XD_NETVAR               // 4
	XD_SAY                  // 5
	XD_MAP                  // 6
	XD_EXITLEVEL            // 7
	XD_ADDFILE              // 8
	XD_PAUSE                // 9
	XD_ADDPLAYER            // 10
	XD_TEAMCHANGE           // 11
	XD_CLEARSCORES          // 12
	XD_LOGIN                // 13
	XD_VERIFIED             // 14
	XD_RANDOMSEED           // 15
	XD_RUNSOC               // 16
	XD_REQADDFILE           // 17
	XD_DELFILE              // 18
	XD_SETMOTD              // 19
	XD_RESPAWN              // 20
	XD_DEMOTED              // 21
	XD_SETUPVOTE            // 22
	XD_MODIFYVOTE           // 23
	XD_PICKVOTE             // 24
	XD_REMOVEPLAYER         // 25
	XD_DISCORD              // 26
	XD_LUACMD               // 27
	XD_LUAVAR               // 28
	MAXNETXCMD
)
