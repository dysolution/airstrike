package airstrike

import "time"

// A Mission is a plan for the execution of "attacks" against an API and a
// representation of its history since inception, including the number of
// attacks that have occurred and how many Planes make up its Squadron.
//
// Enabled missions are eligible for indefinite execution, with each attack
// commencing every Interval seconds.
type Mission struct {
	Enabled   bool      `json:"enabled"`
	Inception time.Time `json:"inception"`
	Interval  int       `json:"interval"`
	RaidCount int       `json:"raid_count"`
	Planes    []Plane   `json:"planes"`
}

// SetInterval changes the length of the pause that occurs between the
// commencement of attacks.
//
// NOTE: Because each Plane deploys its arsenal simultaneously, it is
// possible that a new attack will commence before some of the results from
// the previous attack cycle have been reported. It is possible to overwhelm
// an API's infrastructure if your account is not subject to rate limits.
// The shortest possible interval is one second, because seconds represented
// as integers are easier to reason about than sub-second ones.
func (m *Mission) SetInterval(logCh chan map[string]interface{}, newInterval int) {
	oldInterval := m.Interval
	logCh <- map[string]interface{}{
		"old_interval": oldInterval,
		"new_interval": newInterval,
	}
	m.Interval = newInterval
}
