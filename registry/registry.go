package registry

import (
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/segdumping/micro/registry/conf/etcd"
)

//wrap go-micro registry to support custom config
func NewRegistry(opts ...registry.Option) registry.Registry {
	c, err := etcd.Load()
	if err != nil {
		panic(err)
	}

	switch c.Type {
	case "consul":
		addr := c.Endpoint.Addr + ":" + c.Endpoint.Port
		opts := append(opts, registry.Addrs(addr))
		return consul.NewRegistry(opts...)
	case "etcd":
		addr := c.Endpoint.Addr + ":" + c.Endpoint.Port
		opts := append(opts, registry.Addrs(addr))
		return etcdv3.NewRegistry(opts...)
	}

	return nil
}

