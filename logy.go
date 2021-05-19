package logy

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	AddUUID = true
)

func AdjustAddUUID(enable bool) {
	AddUUID = enable
}

func RegisterWithLogger(logger *logrus.Logger) *logrus.Entry {
	entry := logrus.NewEntry(logger)
	if AddUUID {
		entry = setUUIDToEntry(entry)
	}

	return entry
}

func RegisterTitleWithLogger(funcName string, logger *logrus.Logger) *logrus.Entry {
	entry := logrus.NewEntry(logger)
	if AddUUID {
		entry = setUUIDToEntry(entry)
	}

	return entry.WithField(KeyFuncName, funcName)
}

func RegisterTitleWithEntry(funcName string, entry *logrus.Entry) *logrus.Entry {
	return entry.WithField(KeyFuncName, funcName)
}

func RegisterField(key string, value interface{}, entry *logrus.Entry) *logrus.Entry {
	return entry.WithField(key, value)
}

func GetUUID(entry *logrus.Entry) string {
	u, ok := entry.Data[KeyUUID]
	if ok {
		return fmt.Sprint(u)
	}

	return ""
}

func GetFileName(discardBefore string, dep int) string {
	_, file, _, ok := runtime.Caller(1 + dep)
	if !ok {
		return ""
	}

	if "" == discardBefore {
		return file
	}

	index := strings.Index(file, discardBefore)
	return file[index:]
}
