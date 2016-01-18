package ordnance

import (
	"github.com/Sirupsen/logrus"
	"github.com/dysolution/sleepwalker"
)

type ArmedWeapon interface {
	Fire(sleepwalker.RESTClient) (sleepwalker.Result, error)
	String() string
}

type Arsenal []ArmedWeapon

func NewArsenal(weapons ...ArmedWeapon) Arsenal {
	var arsenal Arsenal
	for _, weapon := range weapons {
		arsenal = append(arsenal, weapon)
	}
	logrus.Debugf("getArsenal: created arsenal: %v", arsenal)
	return arsenal
}
