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

// NewRaid returns a new Raid. It notifies the user if no planes were given.
func NewRaid(inputPlanes ...Plane) (Raid, error) {
	if len(inputPlanes) == 0 {
		return Raid{}, errors.New("no planes to launch")
	}
	var planes []Plane
	for _, plane := range inputPlanes {
		planes = append(planes, plane)
	}
	return Raid{Planes: planes}, nil
}
