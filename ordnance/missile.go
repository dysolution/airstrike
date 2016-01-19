package ordnance

import (
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/sleepwalker"
)

// A Missile represents an action the API client performs whose URL isn't
// known until runtime, such as the retrieval or deletion of the most
// recently created Batch.
type Missile struct {
	Client    sleepwalker.RESTClient
	Name      string                                                   `json:"name"`
	Operation func(sleepwalker.RESTClient) (sleepwalker.Result, error) `json:"-"`
}

// Fire deploys the Missile.
func (m Missile) Fire(c sleepwalker.RESTClient) (sleepwalker.Result, error) {
	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]

	result, err := m.Operation(m.Client)
	if err != nil {
		result.Log().WithFields(logrus.Fields{
			"name":  m.Name,
			"error": err,
		}).Errorf(desc)
		return sleepwalker.Result{}, err
	}
	result.Log().WithFields(logrus.Fields{
		"name": m.Name,
	}).Debug(desc)
	return result, nil
}

func (m Missile) String() string {
	return "Missile: " + m.Name
}
