package airstrike

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/sleepwalker"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var log *logrus.Logger

func init() {
	// log = espsdk.Log
	log = logrus.New()
	log.Formatter = &prefixed.TextFormatter{TimestampFormat: time.RFC3339}
}

type ArmedWeapon interface {
	Fire(sleepwalker.RESTClient) (sleepwalker.Result, error)
}

// A Plane has an arsenal of deployable weapons. It represents a list of
// tasks that, perfored serially, compose a workflow.
//
// Many planes can deploy their arsenal at the same time, but each
// weapon in a plane's arsenal must be deployed one at a time.
//
// For example, a common workflow would be:
//   1. list all batches
//   2. get the metadata for a batch
//   3. upload a contribution to the batch
type Plane struct {
	Name    string `json:"name"`
	Client  sleepwalker.RESTClient
	Arsenal []ArmedWeapon `json:"weapons"`
}

// Deploy sequentially fires all of the weapons within an Arsenal and reports
// the results.
func Deploy(p Plane) ([]sleepwalker.Result, error) {
	var results []sleepwalker.Result
	for _, weapon := range p.Arsenal {
		log.Debugf("deploying %s", weapon)
		result, _ := weapon.Fire(p.Client)
		results = append(results, result)
	}
	return results, nil
}
