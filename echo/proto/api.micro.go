// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: api.proto

package go_api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	proto1 "github.com/micro/go-micro/api/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Echo service

type EchoService interface {
	Say(ctx context.Context, in *proto1.Request, opts ...client.CallOption) (*proto1.Response, error)
	Hi(ctx context.Context, in *proto1.Request, opts ...client.CallOption) (*proto1.Response, error)
}

type echoService struct {
	c    client.Client
	name string
}

func NewEchoService(name string, c client.Client) EchoService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.api"
	}
	return &echoService{
		c:    c,
		name: name,
	}
}

func (c *echoService) Say(ctx context.Context, in *proto1.Request, opts ...client.CallOption) (*proto1.Response, error) {
	req := c.c.NewRequest(c.name, "Echo.Say", in)
	out := new(proto1.Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *echoService) Hi(ctx context.Context, in *proto1.Request, opts ...client.CallOption) (*proto1.Response, error) {
	req := c.c.NewRequest(c.name, "Echo.Hi", in)
	out := new(proto1.Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Echo service

type EchoHandler interface {
	Say(context.Context, *proto1.Request, *proto1.Response) error
	Hi(context.Context, *proto1.Request, *proto1.Response) error
}

func RegisterEchoHandler(s server.Server, hdlr EchoHandler, opts ...server.HandlerOption) error {
	type echo interface {
		Say(ctx context.Context, in *proto1.Request, out *proto1.Response) error
		Hi(ctx context.Context, in *proto1.Request, out *proto1.Response) error
	}
	type Echo struct {
		echo
	}
	h := &echoHandler{hdlr}
	return s.Handle(s.NewHandler(&Echo{h}, opts...))
}

type echoHandler struct {
	EchoHandler
}

func (h *echoHandler) Say(ctx context.Context, in *proto1.Request, out *proto1.Response) error {
	return h.EchoHandler.Say(ctx, in, out)
}

func (h *echoHandler) Hi(ctx context.Context, in *proto1.Request, out *proto1.Response) error {
	return h.EchoHandler.Hi(ctx, in, out)
}