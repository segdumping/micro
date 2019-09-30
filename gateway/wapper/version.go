package wapper

import (
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/client/selector"
	metadata2 "github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
)

type versionWrapper struct {
	client.Client
}

//version wrapper
func (v *versionWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	filter := func(services []*registry.Service) []*registry.Service {
		metadata, ok := metadata2.FromContext(ctx)
		if !ok {
			return services
		}

		version, ok := metadata["Version"]
		if !ok {
			return services
		}

		var valids []*registry.Service
		for _, service := range services {
			if service.Version == version {
				valids = append(valids, service)
			}
		}

		return valids
	}


	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(filter),
		))

	return v.Client.Call(ctx, req, rsp, callOptions...)
}

func NewVersionWrapper() client.Wrapper {
	return func (c client.Client) client.Client {
		return &versionWrapper{c}
	}
}
