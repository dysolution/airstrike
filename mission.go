package airstrike

import "time"

type Mission struct {
	Enabled      bool
	Inception    time.Time
	Interval     int
	RaidCount    int
	SquadronSize int
}

func (m *Mission) SetInterval(interval int) {
}
