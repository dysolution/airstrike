package airstrike

import (
	"encoding/json"
	"sync"

	"github.com/dysolution/sleepwalker"
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

// Conduct concurrently drops all of the Bombs in a Raid's Payload and
// returns a collection of the results.
func (r *Raid) Conduct() ([]sleepwalker.Result, error) {
	var allResults []sleepwalker.Result
	var reporterWg = sync.WaitGroup{}
	var ch chan sleepwalker.Result

	for arsenalID, arsenal := range r.Planes {
		squadron := NewSquadron()
		go squadron.Bombard(ch, arsenalID, arsenal)
		go func() {
			reporterWg.Add(1)
			result := <-ch
			allResults = append(allResults, result)
		}()
	}
	return allResults, nil
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
