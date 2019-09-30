package conf

import (
	"encoding/xml"
	"errors"
	"github.com/segdumping/micro/log"
	"io/ioutil"
	"path"
	"runtime"
)

var Config config

type config struct {
	Micro micro `xml:"micro"`
	Log log.Config `xml:"log"`
}

type micro struct {
	CheckDuration int `xml:"checkDuration"`
	ReapDuration  int `xml:"reapDuration"`
}

func Load() error {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("runtime caller error")
	}

	b, err := ioutil.ReadFile(path.Dir(f) + "/config.xml")
	if err != nil {
		return err
	}

	return xml.Unmarshal(b, &Config)
}
