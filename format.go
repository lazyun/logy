package logy

import (
	"github.com/sirupsen/logrus"
)

type format struct{}

func (f format) Format(entry *logrus.Entry) ([]byte, error) {
	//msg := fmt.Sprintf("")
	//return []byte()
	return nil, nil
}
