package local

import (
	"encoding/xml"
	"errors"
	"github.com/segdumping/micro/registry/conf"
	"io/ioutil"
	"path"
	"runtime"
)

func Load() (*conf.Registry, error) {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("runtime caller error")
	}

	b, err := ioutil.ReadFile(path.Dir(f) + "/config.xml")
	if err != nil {
		return nil, err
	}

	var config conf.Config
	err = xml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return config.GetValid()
}
