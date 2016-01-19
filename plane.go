package airstrike

import (
	"errors"
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike/ordnance"
	"github.com/dysolution/sleepwalker"
)

var log = logrus.New()

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
	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]
	log.WithFields(logrus.Fields{
		"name":   name,
		"client": client,
	}).Debug(desc)
	return Plane{Name: name, Client: client}
}

// Arm loads the given arsenal into the Plane and logs error conditions.
func (p *Plane) Arm(weapons ordnance.Arsenal) {
	desc := "airstrike.(*Plane).Arm"
	if len(weapons) == 0 {
		log.WithFields(logrus.Fields{
			"plane":   p.Name,
			"weapons": weapons,
			"error":   errors.New("no weapons provided"),
		}).Error(desc)
	} else {
		log.WithFields(logrus.Fields{
			"plane":   p.Name,
			"weapons": weapons,
		}).Debug(desc)
		p.Arsenal = weapons
	}
}

// Launch tells a Plane to sequentially fires all of its weapons and report
// the results.
func (p Plane) Launch() ([]sleepwalker.Result, error) {
	var results []sleepwalker.Result

	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]

	log.WithFields(logrus.Fields{
		"plane": p,
	}).Info(desc)

	for _, weapon := range p.Arsenal {

		if weapon == nil {
			log.WithFields(logrus.Fields{
				"error":  "nil weapon",
				"plane":  p,
				"weapon": weapon,
			}).Error(desc)
			continue
		}

		log.WithFields(logrus.Fields{
			"client": p.Client,
			"msg":    "firing weapon",
			"plane":  p,
			"weapon": weapon,
		}).Debug(desc)

		if p.Client == nil {
			log.WithFields(logrus.Fields{
				"plane": p,
				"error": "nil client",
			}).Warn(desc)
			continue
		}

		result, _ := weapon.Fire(p.Client) // Fire does its own error logging
		results = append(results, result)
	}
	log.Debugf(desc+" is returning %v results", len(results))
	return results, nil
}
