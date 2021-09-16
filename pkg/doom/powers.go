package doom

type Power = int

const (
	PWInvulnerability Power = iota
	PWSneakers
	PWFlashing
	PWShield
	PWTailsFly   // tails flying
	PWUnderwater // underwater timer
	PWSpaceTime  // In space, no one can hear you spin!
	PWExtraLife  // Extra Life timer

	PWSuper        // Are you super?
	PWGravityBoots // gravity boots

	// Weapon ammunition
	PWInfinityRing
	PWAutomaticRing
	PWBounceRing
	PWScatterRing
	PWGrenadeRing
	PWExplosionRing
	PWRailRing

	// Power Stones
	PWEmeralds // stored like global 'emeralds' variable

	// NiGHTS powerups
	PWNightsSuperloop
	PWNightsHelper
	PWNightsLinkFreeze

	//for linedef exec 427
	PWNoControl
	PWInGoop

	NUMPOWERS
)
