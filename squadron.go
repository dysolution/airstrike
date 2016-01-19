package airstrike

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike/ordnance"
	"github.com/dysolution/sleepwalker"
	"github.com/speps/go-hashids"
)

// A Squadron is a collection of Planes that will simultaneously begin
// deploying their weapons.
type Squadron struct {
	ID     string  `json:"id"`
	Planes []Plane `json:"planes"`
}

// NewSquadron assigns each new Squadron a unique ID and logs its creation.
func NewSquadron(logger *logrus.Logger) Squadron {
	log = logger
	hd := hashids.NewData()
	hd.Salt = "awakened salt"
	hd.MinLength = 4
	h := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{8, 2, 5, int(time.Now().UnixNano())})

	log.WithFields(logrus.Fields{
		"id": id,
	}).Debug("airstrike.NewSquadron")

	return Squadron{id, []Plane{}}
}

// Add associates the provided Plane with the Squadron, logging its addition.
func (s *Squadron) Add(plane Plane) {
	s.Planes = append(s.Planes, plane)
	log.WithFields(logrus.Fields{
		"plane": plane.Name,
	}).Debug("airstrike.(*Squadron).Add")
}

// AddClones creates the specified number of Planes, each armed with the
// payload described, which can be one or more weapons.
func (s *Squadron) AddClones(clones int, client sleepwalker.RESTClient, armory ordnance.Armory, weaponNames ...string) {
	for i := 1; i <= clones; i++ {
		var arsenal ordnance.Arsenal
		name := fmt.Sprintf("clone_%d_of_%d", i, clones)
		for _, weapon := range weaponNames {
			name = fmt.Sprintf("%s_%s", name, weapon)
			arsenal = append(arsenal, armory.GetWeapon(weapon))
		}
		plane := NewPlane(name, client)
		plane.Arm(arsenal)
		s.Add(plane)
	}
}

// AddChaos creates the specified number of Planes that each have their own
// random selection of n weapons, quantified by their "deadliness."
func (s *Squadron) AddChaos(clones int, deadliness int, client sleepwalker.RESTClient, armory ordnance.Armory) {
	desc := "AddChaos"
	for i := 1; i <= clones; i++ {
		name := fmt.Sprintf("chaos_%d_of_%d", i, clones)
		plane := NewPlane(name, client)

		weaponNames := armory.GetRandomWeaponNames(deadliness)
		var arsenal ordnance.Arsenal
		for _, weaponName := range weaponNames {
			log.WithFields(logrus.Fields{
				"plane":  name,
				"weapon": weaponName,
				"msg":    "adding weapon",
			}).Debug(desc)
			arsenal = append(arsenal, armory.GetWeapon(weaponName))
		}
		plane.Arm(arsenal)
		s.Add(plane)
	}
}
