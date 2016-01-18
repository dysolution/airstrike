package ordnance

import (
	"errors"

	"github.com/dysolution/sleepwalker"
)

import "github.com/Sirupsen/logrus"

var log *logrus.Logger

func init() {
}

type Armory struct {
	Weapons map[string]ArmedWeapon
}

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
