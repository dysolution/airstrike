package ordnance

import "github.com/dysolution/sleepwalker"

// A Missile represents an action the API client performs whose URL isn't
// known until runtime, such as the retrieval or deletion of the most
// recently created Batch.
type Missile struct {
	Client    sleepwalker.RESTClient
	Name      string                                                   `json:"name"`
	Operation func(sleepwalker.RESTClient) (sleepwalker.Result, error) `json:"-"`
}

// Fire deploys the Missile.
func (m Missile) Fire(c sleepwalker.RESTClient, logCh chan map[string]interface{}) (sleepwalker.Result, error) {
	desc := "airstrike.ordnance.Missile.Fire"

	result, err := m.Operation(m.Client)
	if err != nil {
		result.Log().WithFields(map[string]interface{}{
			"name":  m.Name,
			"error": err,
		}).Errorf(desc)
		return sleepwalker.Result{}, err
	}
	logCh <- map[string]interface{}{
		"name":   m.Name,
		"source": desc,
	}
	return result, nil
}

func (m Missile) String() string {
	return "Missile: " + m.Name
}
