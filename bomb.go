package airstrike

import (
	"errors"
	"fmt"

	"github.com/dysolution/espsdk"
)

// A Bomb represents an action for the API client to perform. Any API
// operation that doesn't require knowledge of the state of your account can
// be a Bomb.
//
// If the API client will need to inspect your account before performing the
// action, e.g., delete the most-recently-created Submission Batch, you should
// use a Missile instead.
type Bomb struct {
	Client  *espsdk.Client
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Payload espsdk.RESTObject `json:"payload,omitempty"`
}

// String implements Stringer.
func (b *Bomb) String() string {
	return "Bomb: " + b.Name
}

// Fire deploys the Bullet.
func (b Bomb) Fire() (espsdk.Result, error) {
	switch b.Method {
	case "GET", "get":
		return b.handler(b.Client.Get)
	case "POST", "post":
		return b.handler(b.Client.Create)
	case "PUT", "put":
		return b.handler(b.Client.Update)
	case "DELETE", "delete":
		return b.handler(b.Client.Delete)
	}
	msg := fmt.Sprintf("%s.Deploy: undefined method: %s", b.Name, b.Method)
	return espsdk.Result{}, errors.New(msg)
}

func (b *Bomb) handler(fn func(espsdk.Findable) (espsdk.Result, error)) (espsdk.Result, error) {
	result, err := fn(b.Payload)
	if err != nil {
		log.Errorf("%s.Deploy %s: %v", b.Name, b.Method, err)
		return espsdk.Result{}, err
	}
	result.Log().Debugf("%s.Deploy", b.Name)
	return result, nil
}
