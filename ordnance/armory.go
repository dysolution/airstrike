package ordnance

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/dysolution/sleepwalker"
)

import "github.com/Sirupsen/logrus"

var log *logrus.Logger

func init() {
}

// An Armory maintains a collection of weapons that can be retrieved by name
// or at random.
type Armory struct {
	Weapons map[string]ArmedWeapon `json:"weapons"`
}

// NewArmory allows the logger to be specified.
func NewArmory(logger *logrus.Logger) Armory {
	log = logger
	weapons := make(map[string]ArmedWeapon)
	return Armory{Weapons: weapons}
}

func (a *Armory) NewBomb(client sleepwalker.RESTClient, name string, method string, url string, payload sleepwalker.RESTObject) {
	var payloadPresent bool
	if payload != nil {
		payloadPresent = true
	} else {
		payloadPresent = false
	}
	log.WithFields(logrus.Fields{
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
	log.WithFields(logrus.Fields{
		"name":      name,
		"operation": op,
	}).Debugf("Armory.NewMissile")
	a.Weapons[name] = Missile{
		Client:    client,
		Name:      name,
		Operation: op,
	}
}

func (a Armory) GetWeapon(name string) ArmedWeapon {
	desc := "Armory.GetWeapon"
	if a.Weapons[name] == nil {
		err := errors.New("undefined weapon")
		log.WithFields(logrus.Fields{
			"name":  name,
			"error": err,
		}).Error(desc)
		return nil
	}
	log.WithFields(logrus.Fields{
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

	log.WithFields(logrus.Fields{
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
