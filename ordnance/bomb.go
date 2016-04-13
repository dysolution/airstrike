package ordnance

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dysolution/sleepwalker"
)

// A Bomb represents an action for the API client to perform. Any API
// operation that doesn't require knowledge of the state of your account can
// be a Bomb.
//
// If the API client will need to inspect your account before performing
// the action, e.g., delete the most-recently-created Submission Batch,
// you should use a Missile instead.
type Bomb struct {
	Client  sleepwalker.RESTClient
	Name    string                 `json:"name"`
	Method  string                 `json:"method"`
	URL     string                 `json:"url"`
	Payload sleepwalker.RESTObject `json:"payload,omitempty"`
}

// String implements fmt.Stringer.
func (b Bomb) String() string {
	return "Bomb: " + b.Name
}

// Fire deploys the Bomb.
func (b Bomb) Fire(c sleepwalker.RESTClient) (sleepwalker.Result, error) {
	switch strings.ToUpper(b.Method) {
	case "GET":
		return b.Client.Get(b.Payload)
	case "POST":
		return b.Client.Create(b.Payload)
	case "PUT":
		return b.Client.Update(b.Payload)
	case "DELETE":
		return b.Client.Delete(b.Payload)
	}
	msg := fmt.Sprintf("%s.Fire: undefined method: %s", b.Name, b.Method)
	return sleepwalker.Result{}, errors.New(msg)
}
