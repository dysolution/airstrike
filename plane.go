package airstrike

import (
	"errors"
	"strings"
	"time"

	"github.com/dysolution/airstrike/ordnance"
	"github.com/dysolution/sleepwalker"
)

// A Plane has an arsenal of deployable weapons. It represents a list of
// tasks that, perfored serially, compose a workflow.
//
// Many planes can deploy their arsenal at the same time, but each
// weapon in a plane's arsenal must be deployed one at a time.
//
// For example, a common workflow would be:
//   1. GET index - list all items in a collection
//   2. GET show - get the metadata for an item
//   3. POST create - create and/or associate an item related to the first
//
type Plane struct {
	Name    string                 `json:"name"`
	Client  sleepwalker.RESTClient `json:"-"`
	Arsenal ordnance.Arsenal       `json:"arsenal"`
}

// NewPlane ensures that the creation of each Plane is logged.
func NewPlane(name string, client sleepwalker.RESTClient) Plane {
	desc := "airstrike.NewPlane"
	log.WithFields(map[string]interface{}{
		"name":   name,
		"client": client,
	}).Debug(desc)
	return Plane{Name: name, Client: client}
}

// Arm loads the given arsenal into the Plane and logs error conditions.
func (p *Plane) Arm(weapons ordnance.Arsenal) {
	desc := "airstrike.(*Plane).Arm"
	if len(weapons) == 0 {
		log.WithFields(map[string]interface{}{
			"plane":   p.Name,
			"weapons": weapons,
			"error":   errors.New("no weapons provided"),
		}).Error(desc)
	} else {
		log.WithFields(map[string]interface{}{
			"plane":   p.Name,
			"weapons": weapons,
		}).Debug(desc)
		p.Arsenal = weapons
	}
}

// Launch tells a Plane to sequentially fires all of its weapons and report
// the results.
func (p Plane) Launch(logCh chan map[string]interface{}) ([]sleepwalker.Result, error) {
	desc := "airstrike.Plane.Launch"
	var results []sleepwalker.Result

	logCh <- map[string]interface{}{
		"source": desc,
		"plane":  p,
	}

	for _, weapon := range p.Arsenal {
		result, err := p.fireWeapon(weapon, logCh)
		if err != nil {
			logCh <- map[string]interface{}{
				"source":   desc,
				"error":    err,
				"plane":    p,
				"weapon":   weapon,
				"severity": "ERROR",
			}
		}
		results = append(results, result)
	}
	return results, nil
}

// runs in a goroutine (Raid.Conduct)
func (p Plane) launchAndReport(urlInvariant string, logCh chan map[string]interface{}, squadronID string) {
	results, err := p.Launch(logCh)
	if err != nil {
		logCh <- map[string]interface{}{
			"error":    err,
			"severity": "ERROR",
		}
	}
	for weaponID, result := range results {
		var path string
		parts := strings.SplitAfter(result.Path, urlInvariant)
		if len(parts) >= 2 {
			path = parts[1]
		} else {
			path = result.Path
		}

		stats := map[string]interface{}{
			"plane":         p.Name,
			"weapon_id":     weaponID,
			"squadron_id":   squadronID,
			"method":        result.Verb,
			"path":          path,
			"response_time": result.Duration * time.Millisecond,
			"status_code":   result.StatusCode,
		}
		logCh <- stats
	}
}
func (p Plane) fireWeapon(weapon ordnance.ArmedWeapon, logCh chan map[string]interface{}) (sleepwalker.Result, error) {
	desc := "airstrike.Plane.fireWeapon"
	if weapon == nil {
		return sleepwalker.Result{}, errors.New("nil weapon")
	}
	if p.Client == nil {
		return sleepwalker.Result{}, errors.New("nil client")
	}

	logCh <- map[string]interface{}{
		"client": p.Client,
		"plane":  p,
		"weapon": weapon,
		"source": desc,
	}

	result, _ := weapon.Fire(p.Client, logCh) // Fire does its own error logging
	return result, nil
}
