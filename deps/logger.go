package deps

import (
	"github.com/sirupsen/logrus"
)

func NewDefaultLogger() *logrus.Logger {
	return logrus.New()
}
