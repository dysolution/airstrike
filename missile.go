package airstrike

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
func (m Missile) Fire(c sleepwalker.RESTClient) (sleepwalker.Result, error) {
	result, err := m.Operation(m.Client)
	if err != nil {
		result.Log().Errorf("%s.Deploy %v: %v", m.Name, m.Operation, err)
		return sleepwalker.Result{}, err
	}
	result.Log().Debugf("%s.Deploy", m.Name)
	return result, nil
}

func (m Missile) String() string {
	return "Missile: " + m.Name
}
