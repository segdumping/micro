package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro"
	api "github.com/micro/go-micro/api/proto"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
	proto "github.com/segdumping/micro/echo/proto"
	"github.com/segdumping/micro/gateway/options"
	"github.com/segdumping/micro/registry"
	"strings"
)

type Echo struct {

}

func (e *Echo) Say(context context.Context, req *api.Request, rsp *api.Response) error {
	log.Logf("api endpoint say")
	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.New("no content")
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode("say to " + strings.Join(name.Values, ","))

	rsp.StatusCode = 200
	rsp.Body = buf.String()

	return nil
}

func (e *Echo) Hi(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Logf("api endpoint hi")
	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.New("no content")
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode("hi " + strings.Join(name.Values, ","))

	rsp.StatusCode = 200
	rsp.Body = buf.String()

	return nil
}

func main() {
	//create server with options
	server := server.NewServer(
		server.Id("1"),
		server.Name("go.micro.api.echo"),
		server.Registry(registry.NewRegistry()),
	)

	service := micro.NewService(
		micro.Server(server),
	)

	service.Init()
	h := new(Echo)

	//add handler type
	options := options.HandlerOptions(h, "api")
	options = append(options)
	_ = proto.RegisterEchoHandler(
		service.Server(),
		h,
		options...
		)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
