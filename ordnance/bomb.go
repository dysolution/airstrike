package ordnance

import (
	"errors"
	"fmt"

	"github.com/dysolution/sleepwalker"
)

// A Bomb represents an action for the API client to perform. Any API
// operation that doesn't require knowledge of the state of your account can
// be a Bomb.
//
// If the API client will need to inspect your account before performing the
// action, e.g., delete the most-recently-created Submission Batch, you should
// use a Missile instead.
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

// func (b Bomb) Client() sleepwalker.RESTClient {
// 	return b.Client
// }

// Fire deploys the Bullet.
func (b Bomb) Fire(c sleepwalker.RESTClient) (sleepwalker.Result, error) {
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
	msg := fmt.Sprintf("%s.Fire: undefined method: %s", b.Name, b.Method)
	return sleepwalker.Result{}, errors.New(msg)
}

func (b Bomb) NoPayloadError(desc string) error {
	msg := fmt.Sprintf("%v: payload for %v: %v", desc, b, b.Payload)
	return errors.New(msg)
}

func (b Bomb) handler(fn func(sleepwalker.Findable) (sleepwalker.Result, error)) (sleepwalker.Result, error) {
	result, err := fn(b.Payload)
	if err != nil {
		return sleepwalker.Result{}, err
	}
	return result, nil
}

func (b Bomb) log(logCh chan map[string]interface{}, severity, desc string, err error) {
	logCh <- map[string]interface{}{
		"name":     b.Name,
		"method":   b.Method,
		"path":     b.URL,
		"error":    err,
		"source":   desc,
		"severity": severity,
	}
}
