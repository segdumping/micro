package wapper

import (
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/sirupsen/logrus"
)

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	logrus.Infof("ctx: %v service: %s method: %s", md, req.Service(), req.Endpoint())
	return l.Client.Call(ctx, req, rsp, opts...)
}

type traceWrapper struct {
	client.Client
}

func (t *traceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, ok := metadata.FromContext(ctx)
	if ok {
		md["trace-id"] = uuid.New().String()
		ctx = metadata.NewContext(ctx, md)
	}

	return t.Client.Call(ctx, req, rsp, opts...)
}

func LogWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

func TraceWrap(c client.Client) client.Client {
	return &traceWrapper{c}
}
