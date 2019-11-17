package tools

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
)

func StructToIOBody(i interface{}) io.Reader {
	b, err := json.Marshal(i)
	if err != nil {
		logrus.Errorln(err)
	}
	return bytes.NewReader(b)

}
