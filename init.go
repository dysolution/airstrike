package airstrike

import "github.com/Sirupsen/logrus"

var log *logrus.Logger

func init() {
	log = logrus.New()
}
