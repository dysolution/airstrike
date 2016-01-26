package ordnance

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/dysolution/sleepwalker"
)

import "github.com/Sirupsen/logrus"

var log *logrus.Logger

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
func NewArmory(logger *logrus.Logger) Armory {
	log = logger
	weapons := make(map[string]Weapon)
	return Armory{Weapons: weapons}
}

func (a *Armory) NewBomb(client sleepwalker.RESTClient, name string, method string, url string, payload sleepwalker.RESTObject) {
	var payloadPresent bool
	if payload != nil {
		payloadPresent = true
	} else {
		payloadPresent = false
	}
	log.WithFields(map[string]interface{}{
		"name":             name,
		"method":           method,
		"path":             url,
		"payload_included": payloadPresent,
	}).Debugf("Armory.NewBomb")

	a.Weapons[name] = Bomb{
		Client:  client,
		Name:    name,
		Method:  method,
		URL:     url,
		Payload: payload,
	}
}

func (a *Armory) NewMissile(client sleepwalker.RESTClient, name string, op func(sleepwalker.RESTClient) (sleepwalker.Result, error)) {
	log.WithFields(map[string]interface{}{
		"name": name,
	}).Debugf("Armory.NewMissile")
	a.Weapons[name] = Missile{
		Client:    client,
		Name:      name,
		Operation: op,
	}
}

func (a Armory) GetWeapon(name string) Weapon {
	desc := "Armory.GetWeapon"
	if a.Weapons[name] == nil {
		err := errors.New("undefined weapon")
		log.WithFields(map[string]interface{}{
			"name":  name,
			"error": err,
		}).Error(desc)
		return nil
	}
	log.WithFields(map[string]interface{}{
		"name": name,
	}).Debug(desc)
	return a.Weapons[name]
}

func (a *Armory) GetRandomWeaponNames(count int) []string {
	desc := "Armory.GetRandomWeaponNames"
	var names []string
	for i := count; i > 0; i-- {
		names = append(names, a.getRandomWeapon())
	}

	log.WithFields(map[string]interface{}{
		"weapons": names,
	}).Debug(desc)
	return names
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
