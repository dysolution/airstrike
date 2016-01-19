package airstrike

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dysolution/airstrike/ordnance"
	"github.com/dysolution/sleepwalker"
	"github.com/speps/go-hashids"
)

type Squadron struct {
	wg     sync.WaitGroup
	ID     string  `json:"id"`
	Planes []Plane `json:"planes"`
}

func New(logger *logrus.Logger) Squadron {
	log = logger
	log.Debugf("creating squadron")
	var wg sync.WaitGroup

	hd := hashids.NewData()
	hd.Salt = "awakened salt"
	hd.MinLength = 4
	h := hashids.NewWithData(hd)
	id, _ := h.Encode([]int{8, 2, 5, int(time.Now().UnixNano())})

	return Squadron{wg, id, []Plane{}}
}

func (s *Squadron) Add(plane Plane) {
	myPC, _, _, _ := runtime.Caller(0)
	desc := runtime.FuncForPC(myPC).Name()
	desc = strings.SplitAfter(desc, "github.com/dysolution/")[1]

	s.Planes = append(s.Planes, plane)

	log.WithFields(logrus.Fields{
		"plane":  plane.Name,
		"planes": s.Planes,
	}).Debug(desc)
}

// AddClones creates the specified number of Planes, each armed with the
// payload described, which can be one or more weapons.
func (s *Squadron) AddClones(clones int, client sleepwalker.RESTClient, armory ordnance.Armory, weapons ...string) {
	for i := 1; i <= clones; i++ {
		name := fmt.Sprintf("clone_%d_of_%d", i, clones)
		for _, weapon := range weapons {
			name = fmt.Sprintf("%s_%s", name, weapon)
		}
		plane := NewPlane(name, client)
		for _, weapon := range weapons {
			plane.Arm(armory.GetArsenal(weapon))
		}
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
