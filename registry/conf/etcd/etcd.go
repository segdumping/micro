package etcd

import (
	"encoding/json"
	"errors"
	"github.com/segdumping/micro/configcenter"
	"github.com/segdumping/micro/registry/conf"
)

const (
	key = "custom-registry"
)

func Load() (*conf.Registry, error) {
	m ,err := configcenter.Get(key)
	if err != nil {
		return nil, err
	}

	//check expect
	r, ok := m[key]
	if !ok {
		return nil, errors.New("remote empty")
	}

	var config conf.Config
	err = json.Unmarshal(r, &config)
	if err != nil {
		return nil, err
	}

	return config.GetValid()
}
