package dependencies

import (
	"os"

	"github.com/sirupsen/logrus"
)

func GetLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.InfoLevel)
	return log
}
