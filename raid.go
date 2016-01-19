package airstrike

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

// A SimpleRaid reports only the name of each object to simplify output.
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

func writeLog(responseTime, warningThreshold time.Duration, fields logrus.Fields) {
	desc := "airstrike.report"
	if responseTime > warningThreshold {
		log.WithFields(fields).Warn(desc)
	} else {
		log.WithFields(fields).Info(desc)
	}
}

func charsToSave(responseTime time.Duration) int {
	if responseTime/time.Millisecond >= 1000 {
		return 6
	}
	return 5
}

func writeConsoleLegend() {
	for i := 0; i <= 8; i++ {
		for j := 0; j < 10; j++ {
			if j == 0 {
				fmt.Printf("%d", i)
			} else {
				fmt.Print(" ")
			}

		}
	}
	fmt.Printf("\r")
}

func writeConsoleGauge(responseTime time.Duration) {
	// 80 chars * 1 block per 10 ms = 800 ms max resolution
	maxRes := 800
	maxWidth := 80
	blockMS := time.Duration(maxRes) / time.Duration(maxWidth)
	numBlocks := int(responseTime / blockMS / time.Millisecond)
	// fmt.Print("[")
	writeConsoleLegend()

	// allow chars for "nnnnms" text
	for i := 0; i <= numBlocks-charsToSave(responseTime); i++ {
		if i <= maxWidth-charsToSave(responseTime) {
			fmt.Print("#")
		}
	}
	fmt.Printf("%dms", responseTime/time.Millisecond)
	// for i := 1; i < charsToSave(responseTime); i++ {
	// 	fmt.Print("-")
	// }
	fmt.Println()
}

func report(ch chan logrus.Fields, urlInvariant string, warningThreshold time.Duration) {
	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]
	for {
		fields := <-ch
		responseTime, _ := fields["response_time"].(time.Duration)
		writeLog(responseTime, warningThreshold, fields)
		writeConsoleGauge(responseTime)
	}
}

// Conduct tells the Squadron to launch all of its planes. Each Plane serially
// fires its weapons and sends the result of each weapon down a channel.
func (r *Raid) Conduct(logger *logrus.Logger, urlInvariant string, warningThreshold time.Duration) {
	ch := make(chan logrus.Fields)

	squadron := NewSquadron(logger)

	go report(ch, urlInvariant, warningThreshold)

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
