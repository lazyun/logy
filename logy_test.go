package logy

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestRegister(t *testing.T) {
	log := logrus.New()
	e := RegisterTitleWithLogger("TestRegister", log)
	e = RegisterTitleWithEntry("TestRegister1", e)

	e.Info("biu~biu~")
}

func TestGetFileName(t *testing.T) {
	t.Log(GetFileName("logy"))
}
