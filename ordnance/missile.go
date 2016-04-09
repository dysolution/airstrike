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
func (m Missile) Fire(c sleepwalker.RESTClient) (result sleepwalker.Result, err error) {
	result, err = m.Operation(m.Client)
	return
}

func (m Missile) String() string {
	return "Missile: " + m.Name
}
