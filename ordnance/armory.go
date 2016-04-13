// Package ordnance is concerned with constructing and storing weapons and
// providing access to them by name.
package ordnance

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	sw "github.com/dysolution/sleepwalker"
)

// var log *logrus.Logger

func init() {
}

type NoPayloadError struct {
	msg string
	obj interface{}
}

func (e NoPayloadError) Error() string {
	return fmt.Sprintf("%s: nil payload: %v", e.msg, e.obj)
}

// An Armory maintains a collection of weapons that can be retrieved by name
// or at random.
type Armory struct {
	Weapons map[string]Weapon `json:"weapons"`
}

// NewArmory allows the logger to be specified.
func NewArmory() Armory {
	weapons := make(map[string]Weapon)
	return Armory{Weapons: weapons}
}

func (a *Armory) NewBomb(client sw.RESTClient, name string, method string, url string, payload sw.RESTObject) {
	a.Weapons[name] = Bomb{
		Client:  client,
		Name:    name,
		Method:  method,
		URL:     url,
		Payload: payload,
	}
}

func (a *Armory) NewMissile(client sw.RESTClient, name string, op func(sw.RESTClient) (sw.Result, error)) {
	a.Weapons[name] = Missile{
		Client:    client,
		Name:      name,
		Operation: op,
	}
}

func (a Armory) GetWeapon(name string) Weapon {
	if a.Weapons[name] == nil {
		return nil
	}
	return a.Weapons[name]
}

func (a *Armory) GetRandomWeaponNames(count int) (names []string) {
	for i := count; i > 0; i-- {
		names = append(names, a.getRandomWeapon())
	}
	return
}

func (a *Armory) getRandomWeapon() string {
	rand.Seed(time.Now().UTC().UnixNano())
	m := sync.Mutex{}

	// build a slice of weapon names
	var weapons []string
	m.Lock()
	for name, _ := range a.Weapons {
		weapons = append(weapons, name)
	}
	m.Unlock()

	return weapons[rand.Intn(len(weapons))]
}

func (a Armory) GetArsenal(names ...string) Arsenal {
	var arsenal Arsenal
	for _, name := range names {
		arsenal = append(arsenal, a.GetWeapon(name))
	}
	return arsenal
}
