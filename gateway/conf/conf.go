package conf

import (
	"encoding/xml"
	"errors"
	"github.com/segdumping/micro/log"
	"github.com/segdumping/micro/server/conf"
	"io/ioutil"
	"path"
	"runtime"
)

var Config config

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

type config struct {
	Micro  micro       `xml:"micro"`
	Log    log.Config  `xml:"log"`
	Server conf.Config `xml:"server"`
}

type micro struct {
	Handler   string `xml:"handler"`
	Resolver  string `xml:"resolver"`
	APIPath   string `xml:"apiPath"`
}
