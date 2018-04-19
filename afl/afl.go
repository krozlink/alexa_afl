package afl

import (
	"time"
)

type AFLMatch struct {
	HomeTeam   string
	AwayTeam   string
	HomeOdds   float32
	AwayOdds   float32
	MatchStart time.Time
}

type AFLTeam struct {
	Name      string
	Nicknames []string
}
