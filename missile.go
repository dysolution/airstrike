package airstrike

import "github.com/dysolution/espsdk"

// A Missile represents an action the API client performs whose URL isn't
// known until runtime, such as the retrieval or deletion of the most
// recently created Batch.
type Missile struct {
	Client    *espsdk.Client
	Name      string                        `json:"name"`
	Operation func() (espsdk.Result, error) `json:"-"`
}

// Fire deploys the Missile.
func (m Missile) Fire() (espsdk.Result, error) {
	result, err := m.Operation()
	if err != nil {
		result.Log().Errorf("%s.Deploy %v: %v", m.Name, m.Operation, err)
		return espsdk.Result{}, err
	}
	result.Log().Debugf("%s.Deploy", m.Name)
	return result, nil
}

func (m Missile) String() string {
	return "Missile: " + m.Name
}
