package airstrike

import (
	"encoding/json"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

type SimpleRaid struct {
	Planes []struct {
		Name    string `json:"name"`
		Weapons []struct {
			Name string `json:"name"`
		} `json:"weapons"`
	} `json:"planes"`
}

// A Raid is a collection of bombs capable of reporting summary statistics.
type Raid struct {
	Planes []Plane `json:"planes"`
}

func reportResults(ch chan logrus.Fields, urlInvariant string, warningThreshold time.Duration) {
	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]
	for {
		fields := <-ch
		responseTime, _ := fields["response_time"].(time.Duration)
		if responseTime > warningThreshold {
			log.WithFields(fields).Warn(desc)
		} else {
			log.WithFields(fields).Info(desc)
		}
	}
}

// Conduct tells the Squadron to launch all of its planes. Each Plane serially
// fires its weapons and sends the result of each weapon down a channel.
func (r *Raid) Conduct(logger *logrus.Logger, urlInvariant string, warningThreshold time.Duration) {
	ch := make(chan logrus.Fields)

	squadron := New(logger)

	go reportResults(ch, urlInvariant, warningThreshold)

	for _, plane := range r.Planes {

		go func(plane Plane) {
			myPC, _, _, _ := runtime.Caller(0)
			desc := runtime.FuncForPC(myPC).Name()
			desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]

			results, err := plane.Launch()
			if err != nil {
				ch <- logrus.Fields{"error": err}
			}

			for weaponID, result := range results {
				stats := logrus.Fields{
					"plane":         plane.Name,
					"weapon_id":     weaponID,
					"squadron_id":   squadron.ID,
					"method":        result.Verb,
					"path":          strings.SplitAfter(result.Path, urlInvariant)[1],
					"response_time": result.Duration * time.Millisecond,
					"status_code":   result.StatusCode,
				}
				ch <- stats
			}
		}(plane)

	}
}

func (r *Raid) String() string {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "error marshaling Raid"
	}
	return string(out)
}

// NewRaid initializes and returns a Raid, . It should be used in lieu of Raid literals.
func NewRaid(planes ...Plane) Raid {
	var payload []Plane
	for _, plane := range planes {
		payload = append(payload, plane)
	}
	return Raid{Planes: payload}
}
