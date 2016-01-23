package airstrike

import (
	"encoding/json"
	"errors"

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

// Conduct tells the Squadron to launch all of its planes. Each Plane serially
// fires its weapons and sends the result of each weapon down a channel.
func (r *Raid) Conduct(logger *logrus.Logger, urlInvariant string, logCh chan map[string]interface{}) {
	logCh <- map[string]interface{}{
		"msg":    "creating a squadron",
		"source": "airstrike.Raid.Conduct",
	}
	squadron := NewSquadron(logCh)
	for _, plane := range r.Planes {
		go plane.launchAndReport(urlInvariant, logCh, squadron.ID)
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
