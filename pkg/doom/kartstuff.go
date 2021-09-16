package doom

type KartStuff = int

const (
	// Basic gameplay things
	KPosition      KartStuff = iota // Used for KArt positions, mostly for deterministic stuff
	KOldPosition                    // Used for taunting, when you pass someone
	KPositionDelay                  // Used for position number, so it can grow when passing/being passed
	KPrevCheck                      // Previous checkpoint distance; for p_user.c (was "pw_pcd")
	KNextCheck                      // Next checkpoint distance; for p_user.c (was "pw_ncd")
	KWaypoint                       // Waypoints.
	KStarpostWp                     // Temporarily stores player waypoint for... some reason. Used when respawning and finishing.
	KStarpostFlip                   // the last starpost we hit requires flipping?
	KRespawn                        // Timer for the DEZ laser respawn effect
	KDropDash                       // Charge up for respawn Drop Dash

	KThrowDir      // Held dir of controls; 1 = forward 0 = none -1 = backward (was "player->heldDir")
	KLapAnimation  // Used to show the lap start wing logo animation
	KLapHand       // Lap hand gfx to use; 0 = none 1 = :oKHand: 2 = :thumbs_up: 3 = :thumps_down:
	KCardAnimation // Used to determine the position of some full-screen Battle Mode graphics
	KVoices        // Used to stop the player saying more voices than it should
	KTauntVoices   // Used to specifically stop taunt voice spam
	KInstaShield   // Instashield no-damage animation timer
	KEngineSnd     // Engine sound number you're on.

	KFloorBoost  // Prevents Sneaker sounds for a brief duration, when triggered by a floor panel
	KSpinoutType // Determines whether to thrust forward or not while spinning out; 0 = move forwards 1 = stay still

	KDrift           // Drifting Left or Right plus a bigger counter = sharper turn
	KDriftEnd        // Drift has ended, used to adjust character angle after drift
	KDriftCharge     // Charge your drift so you can release a burst of speed
	KDriftBoost      // Boost you get from drifting
	KBoostCharge     // Charge-up for boosting at the start of the race
	KStartBoost      // Boost you get from start of race or respawn drop dash
	KJmp             // In Mario KArt, letting go of the jump button stops the drift
	KOffroad         // In Super Mario KArt, going offroad has lee-way of about 1 second before you start losing speed
	KPogoSpring      // Pogo spring bounce effect
	KBrakeStop       // Wait until you've made a complete stop for a few tics before letting brake go in reverse.
	KWaterSkip       // Water skipping counter
	KDashpadCooldown // Separate the vanilla SA-style dash pads from using pw_flashing
	KBoostPower      // Base boost value for offroad
	KSpeedBoost      // Boost value smoothing for max speed
	KAccelBoost      // Boost value smoothing for acceleration
	KBoostAngle      // angle set when not spun out OR boosted to determine what direction you should keep going at if you're spun out and boosted.
	KBoostCam        // Camera push forward on boost
	KDestboostCam    // Ditto
	KTimeoverCam     // Camera timer for leaving behind or not
	KAizDriftStrat   // Let go of your drift while boosting? Helper for the SICK STRATZ you have just unlocked
	KBrakeDrift      // Helper for brake-drift spark spawning

	KItemRoulette // Used for the roulette when deciding what item to give you (was "pw_kartitem")
	KRouletteType // Used for the roulette for deciding type (currently only used for Battle to give you better items from KArma items)

	// Item held stuff
	KItemType   // KITEM_ constant for item number
	KItemAmount // Amount of said item
	KItemHeld   // Are you holding an item?

	// Some items use timers for their duration or effects
	//KThunderanim			// Duration of Thunder Shield's use animation
	KCurShield          // 0 = no shield 1 = thunder shield
	KHyudoroTimer       // Duration of the Hyudoro offroad effect itself
	KStealingTimer      // You are stealing an item this is your timer
	KStolenTimer        // You are being stolen from this is your timer
	KSneakerTimer       // Duration of the Sneaker Boost itself
	KGrowshrinkTimer    // > 0 = Big < 0 = small
	KSquishedTimer      // Squished frame timer
	KRocketSneakerTimer // Rocket Sneaker duration timer
	KInvincibilityTimer // Invincibility timer
	KEggmanHeld         // Eggman monitor held separate from KItemheld so it doesn't stop you from getting items
	KEggmanExplode      // Fake item recieved explode in a few seconds
	KEggmanBlame        // Fake item recieved who set this fake
	KLastJawzTarget     // Last person you target with jawz for playing the target switch sfx
	KBananaDrag         // After a second of holding a banana behind you you start to slow down
	KSpinoutTimer       // Spin-out from a banana peel or oil slick (was "pw_bananacam")
	KWipeoutSlow        // timer before you slowdown when getting wiped out
	KJustBumped         // Prevent players from endlessly bumping into each other
	KComebackTimer      // Battle mode how long before you become a bomb after death
	KSadTimer           // How long you've been sad

	// Battle Mode vars
	KBumper         // Number of bumpers left
	KComebackPoints // Number of times you've bombed or gave an item to someone; once it's 3 it gets set back to 0 and you're given a bumper
	KComebackMode   // 0 = bomb 1 = item
	KWanted         // Timer for determining WANTED status lowers when hitting people prevents the game turning into Camp Lazlo
	KYouGotEm       // "You Got Em" gfx when hitting someone as a karma player via a method that gets you back in the game instantly

	// v1.0.2+ vars
	KItemBlink       // Item flashing after roulette prevents Hyudoro stealing AND serves as a mashing indicator
	KItemBlinkMode   // Type of flashing: 0 = white (normal) 1 = red (mashing) 2 = rainbow (enhanced items)
	KGetSparks       // Disable drift sparks at low speed JUST enough to give acceleration the actual headstart above speed
	KJawztargetDelay // Delay for Jawz target switching to make it less twitchy
	KSpectateWait    // How long have you been waiting as a spectator
	KGrowCancel      // Hold the item button down to cancel Grow

	NUMKARTSTUFF
)
