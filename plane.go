package airstrike

import (
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike/arsenal"
	"github.com/dysolution/sleepwalker"
)

var log = logrus.New()

func init() {
	// log = espsdk.Log
	// if log == nil {
	// 	log = logrus.New()
	// 	log.Formatter = &prefixed.TextFormatter{TimestampFormat: time.RFC3339}
	// }
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
	Arsenal arsenal.Arsenal `json:"arsenal"`
}

func NewPlane(name string, client sleepwalker.RESTClient) Plane {
	return Plane{Name: name, Client: client}
}

type EmptyArsenalError struct{}

func (e *EmptyArsenalError) Error() string {
	return "no weapons provided in arsenal"
}

func (p *Plane) Arm(weapons arsenal.Arsenal) error {
	if len(weapons) == 0 {
	}
	p.Arsenal = weapons
	return nil
}

// Deploy sequentially fires all of the weapons within an Arsenal and reports
// the results.
func (p Plane) DeployArsenal() ([]sleepwalker.Result, error) {
	var results []sleepwalker.Result

	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]
	log.Debugf("%s: %v", desc, p)

	for _, weapon := range p.Arsenal {

		log.WithFields(logrus.Fields{
			"client": p.Client,
			"weapon": weapon,
		}).Debug(desc)

		if p.Client == nil {
			log.WithFields(logrus.Fields{
				"plane": p,
			}).Warn(desc)
			continue
		}

		result, err := weapon.Fire(p.Client)
		if err != nil {
			log.Warn(err)
		}
		results = append(results, result)
	}
	return results, nil
}
