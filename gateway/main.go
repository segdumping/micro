package main

import (
	"github.com/gorilla/mux"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/api/resolver"
	rrmicro "github.com/micro/go-micro/api/resolver/micro"
	"github.com/micro/go-micro/api/router"
	regRouter "github.com/micro/go-micro/api/router/registry"
	httpapi "github.com/micro/go-micro/api/server/http"
	"github.com/micro/go-plugins/wrapper/select/roundrobin"
	"github.com/segdumping/micro/gateway/conf"
	"github.com/segdumping/micro/gateway/handler"
	"github.com/segdumping/micro/gateway/wapper"
	"github.com/segdumping/micro/log"
	"github.com/segdumping/micro/registry"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	err := conf.Load()
	if err != nil {
		panic(err)
	}

	err = log.Init(&conf.Config.Log)
	if err != nil {
		panic(err)
	}

	gateway()
}

func gateway() {
	var h http.Handler
	r := mux.NewRouter()
	h = r

	//robin lb
	robin := roundrobin.NewClientWrapper()
	//inner services wrapper
	valid := wapper.NewValidWrapper()
	//version wrapper
	version := wapper.NewVersionWrapper()

	conf := conf.Config

	service := micro.NewService(
		micro.Name(conf.Server.ServerName()),
		micro.WrapClient(wapper.TraceWrap, wapper.LogWrap, valid, version, robin),
		micro.Registry(registry.NewRegistry()),
		micro.RegisterTTL(time.Duration(conf.Server.RegisterTTL) * time.Second),
		micro.RegisterInterval(time.Duration(conf.Server.RegisterInterval) * time.Second),
	)

	// default resolver
	rr := rrmicro.NewResolver(
		resolver.WithNamespace(conf.Server.Namespace),
		resolver.WithHandler(conf.Micro.Handler),
	)

	logrus.Infof("Registering API Request Handler at %s", conf.Micro.APIPath)
	rt := regRouter.NewRouter(
		router.WithNamespace(conf.Server.Namespace),
		router.WithResolver(rr),
		router.WithRegistry(service.Options().Registry),
	)
	r.PathPrefix(conf.Micro.APIPath).Handler(handler.Meta(service, rt))

	// create the server
	api := httpapi.NewServer(conf.Server.Addr)
	api.Init()
	api.Handle("/", h)

	defer func() {
		if err := api.Stop(); err != nil {
			logrus.Fatal(err)
		}
	}()

	if err := api.Start(); err != nil {
		logrus.Fatal(err)
	}

	if err := service.Run(); err != nil {
		logrus.Fatal(err)
	}
}
