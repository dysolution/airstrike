package airstrike

import (
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/sleepwalker"
)

type Squadron struct {
	wg sync.WaitGroup
}

func New(logger *logrus.Logger) Squadron {
	log = logger
	log.Debugf("creating squadron")
	var wg sync.WaitGroup
	return Squadron{wg}
}

func (s Squadron) Bombard(ch chan sleepwalker.Result, pilotID int, plane Plane, squadronID, urlInvariant string) {
	s.wg.Add(1)
	defer s.wg.Done()

	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]
	log.Debugf("%v", desc)

	results, err := plane.DeployArsenal()
	if err != nil {
		log.Errorf("%v: %v", desc, err)
		var result sleepwalker.Result
		ch <- result
	}

	for weaponID, result := range results {
		log.WithFields(logrus.Fields{
			"pilot_id":      pilotID,
			"weapon_id":     weaponID,
			"method":        result.Verb,
			"path":          strings.SplitAfter(result.Path, urlInvariant)[1],
			"response_time": result.Duration * time.Millisecond,
			"status_code":   result.StatusCode,
		}).Info(desc)

		ch <- result
	}
}
