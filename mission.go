package airstrike

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/espsdk"
)

// A Mission is a plan for the execution of "attacks" against an API and a
// representation of its history since inception, including the number of
// attacks that have occurred and how many Planes make up its Squadron.
//
// Enabled missions are eligible for indefinite execution, with each attack
// commencing every Interval seconds.
type Mission struct {
	Enabled         bool          `json:"enabled"`
	EnabledCh       chan bool     `json:"-"`
	Inception       time.Time     `json:"inception"`
	Interval        float64       `json:"interval"`
	IntervalDeltaCh chan float64  `json:"-"`
	MaxResponseTime time.Duration `json:"max_response_time"`
	Planes          []Plane       `json:"planes"`
	RaidCount       int           `json:"raid_count"`
	Reporter        *Reporter     `json:"-"`
	Status          int           `json:"status"`
	StatusCh        chan int      `json:"-"`
}

func NewMission(l *logrus.Logger) *Mission {
	return &Mission{
		Enabled:         false,
		EnabledCh:       make(chan bool, 1),
		Inception:       time.Now(),
		IntervalDeltaCh: make(chan float64, 1),
		Reporter:        NewReporter(l),
	}
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
func (m *Mission) SetInterval(logCh chan map[string]interface{}, newInterval float64) {
	oldInterval := m.Interval
	logCh <- map[string]interface{}{
		"old_interval": oldInterval,
		"new_interval": newInterval,
	}
	m.Interval = newInterval
}

// runs in a goroutine
func (m *Mission) Prosecute(config Raid) {
	desc := "Mission.Prosecute"

	go m.Reporter.Listen()

	m.Reporter.LogCh <- map[string]interface{}{
		"severity": "info",
		"source":   desc,
		"interval": m.Interval,
	}

	for {
		select {
		case m.Enabled = <-m.EnabledCh:
		case m.Status = <-m.StatusCh:
		case d := <-m.IntervalDeltaCh:
			m.setInterval(d)
			m.Reporter.ThresholdCh <- time.Duration(m.Interval) * time.Millisecond
		default:
		}

		if m.Enabled {
			m.Reporter.LogCh <- map[string]interface{}{
				"msg":    "conducting raid",
				"source": desc,
			}
			config.Conduct(espsdk.APIInvariant, m.Reporter.LogCh)
			m.RaidCount++
		}
		m.pauseBetweenRaids()
	}
}

func (m *Mission) pauseBetweenRaids() {
	time.Sleep(time.Duration(m.Interval) * time.Millisecond)
}

func (m *Mission) setInterval(d float64) {
	oldInterval := m.Interval
	m.Interval = float64(m.Interval) + d
	if m.Interval <= 0 {
		m.Interval = float64(0.1)
	}
	m.Reporter.LogCh <- map[string]interface{}{
		"interval_delta": d,
		"interval_new":   m.Interval,
		"interval_old":   oldInterval,
		"message":        "changed interval",
	}
}
