package airstrike

import (
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
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
