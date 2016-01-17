package airstrike

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/sleepwalker"
)

type Squadron struct {
	wg sync.WaitGroup
}

func NewSquadron() Squadron {
	var wg sync.WaitGroup
	return Squadron{wg}
}

func (s *Squadron) Bombard(ch chan sleepwalker.Result, pilotID int, arsenal Plane) {
	s.wg.Add(1)
	defer s.wg.Done()

	results, err := Deploy(arsenal)
	if err != nil {
		log.Errorf("Raid.Conduct(): %v", err)
		var result sleepwalker.Result
		ch <- result
	}

	for weaponID, result := range results {
		log.WithFields(logrus.Fields{
			"pilot_id":  pilotID,
			"weapon_id": weaponID,
			"method":    result.Verb,
			"path":      result.Path,
			// "response_time": result.Duration,
			"response_time": result.Duration * time.Millisecond,
			"status_code":   result.StatusCode,
		}).Info()

		ch <- result
	}
}
