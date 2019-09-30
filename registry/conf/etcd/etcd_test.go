package etcd

import (
	"encoding/json"
	"github.com/segdumping/micro/configcenter"
	"github.com/segdumping/micro/registry/conf"
	"testing"
)

func TestConf(t *testing.T) {
	var c = conf.Config{RegList:[]conf.Registry {
		{
			Type: "consul",
			Valid: false,
			Endpoint:conf.Endpoint{Addr:"127.0.0.1", Port:"8500"},
		},
		{
			Type: "etcd",
			Valid: true,
			Endpoint:conf.Endpoint{Addr:"127.0.0.1", Port:"2379"},
		},
	},
	}

	b, err := json.Marshal(&c)
	if err != nil {
		t.Logf("json marsha1 error: %v", err)
		return
	}


	cc := configcenter.NewConfigCenter()
	_, err = cc.Put("custom-registry", string(b))
	if err != nil {
		t.Logf("config center put error: %v", err)
		return
	}

	t.Log(Load())
}

