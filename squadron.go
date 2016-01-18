package airstrike

import (
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/sleepwalker"
	"github.com/speps/go-hashids"
)

type Squadron struct {
	wg sync.WaitGroup
	ID string
}

func New(logger *logrus.Logger) Squadron {
	log = logger
	log.Debugf("creating squadron")
	var wg sync.WaitGroup

	hd := hashids.NewData()
	hd.Salt = "awakened salt"
	hd.MinLength = 8
	h := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{8, 2, 5, int(time.Now().UnixNano())})

	return Squadron{wg, id}
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
			"squadron_id":   squadronID,
			"method":        result.Verb,
			"path":          strings.SplitAfter(result.Path, urlInvariant)[1],
			"response_time": result.Duration * time.Millisecond,
			"status_code":   result.StatusCode,
		}).Info(desc)

		ch <- result
	}
}
