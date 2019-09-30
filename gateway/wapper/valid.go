package wapper

import (
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/segdumping/micro/configcenter"
	"github.com/segdumping/shared/util"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	key = "invalid-services"
)

type validWrapper struct {
	invalid map[string]bool
	client.Client
}

//this wrapper mean that: a service is registered but not server outer
func (w *validWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	filter := func(services []*registry.Service) []*registry.Service {
		var isDirect bool
		var service string

		metadata, ok := metadata.FromContext(ctx)
		if ok {
			s, ok := metadata["Service"]
			if ok {
				service = s
				isDirect = true
			}
		}

		//specifying service
		if isDirect {
			for _, s := range services {
				for _, node := range s.Nodes {
					if node.Id == service {
						s.Nodes = []*registry.Node{node}
						return []*registry.Service{s}
					}
				}
			}

			return []*registry.Service{}
		}

		//load balance
		for _, service := range services {
			var nodes []*registry.Node
			for _, node := range service.Nodes {
				if !w.invalid[node.Id] {
					nodes = append(nodes, node)
				}
			}
			service.Nodes = nodes
		}
		return services
	}

	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(filter),
	))

	return w.Client.Call(ctx, req, rsp, callOptions...)
}

func (w *validWrapper) Watch(vals map[string][]byte) {
	b, ok := vals[key]
	if !ok {
		return
	}

	l := strings.Split(util.Byte2String(b), ",")
	m := make(map[string]bool, len(l))
	for _, v := range l {
		m[v] = true
	}

	logrus.Debugf("valid watch: %v", m)

	w.invalid = m
}

func NewValidWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		vc := &validWrapper{
			invalid: make(map[string]bool),
			Client:  c,
		}

		vals, err := configcenter.Get(key)
		if err == nil {
			vc.Watch(vals)
		}
		configcenter.Listen(vc, key)

		return vc
	}
}
