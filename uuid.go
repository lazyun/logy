package logy

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type UUIDHook struct{}

func (u UUIDHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.TraceLevel, logrus.DebugLevel,
		logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.FatalLevel}
}

func (u UUIDHook) Fire(entry *logrus.Entry) error {
	_, ok := entry.Data[KeyUUID]
	if ok {
		return nil
	}

	return nil
}

func setUUIDToEntry(entry *logrus.Entry) *logrus.Entry {
	return entry.WithField(KeyUUID, uuid.NewV4().String())
}
