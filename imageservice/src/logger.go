package src

import (
	mgorus "resizeimage/util"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func InitLogger(dbHost, db, collection, username, password string) error {
	logger = logrus.New()
	hooker, err := mgorus.NewHookerWithAuth(dbHost, db, collection, username, password)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logger.Hooks.Add(hooker)
	return nil
}

func LogWithField(object, trace interface{}) *logrus.Entry {
	return logger.WithFields(logrus.Fields{"Object": object, "TraceLog": trace})
}
