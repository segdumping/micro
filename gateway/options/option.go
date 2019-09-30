package options

import (
	"github.com/micro/go-micro/server"
	"reflect"
)

const (
	metaKey = "handler"
)

//struct methods add handler tpye, eg: api、http
func HandlerOptions(v interface{}, handlerType string) (options []server.HandlerOption) {
	if v == nil {
		return
	}

	typ := reflect.TypeOf(v)
	baseName := reflect.Indirect(reflect.ValueOf(v)).Type().Name()
	for i := 0; i < typ.NumMethod(); i++ {
		name := typ.Method(i).Name
		handle := baseName + "." + name
		options = append(options, server.EndpointMetadata(
			handle,
			map[string]string{metaKey: handlerType},
			))
	}

	return
}

//specified method add handler type
func HandlerOption(handler, handlerType string) server.HandlerOption {
	if len(handler) == 0 {
		return emptyOption()
	}

	if len(handlerType) == 0 {
		return emptyOption()
	}

	return server.EndpointMetadata(
		handler,
		map[string]string{metaKey: handlerType},
		)
}

func emptyOption() server.HandlerOption {
	return func(options *server.HandlerOptions) {}
}