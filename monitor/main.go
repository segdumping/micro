package main

import (
	"github.com/micro/go-micro/util/log"
	log2 "github.com/segdumping/micro/log"
	"github.com/segdumping/micro/monitor/conf"
	"github.com/segdumping/micro/monitor/logic"
	"github.com/segdumping/shared/signal"
	"github.com/sirupsen/logrus"

	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	err := conf.Load()
	if err != nil {
		panic(err)
	}

	err = log2.Init(&conf.Config.Log)
	if err != nil {
		panic(err)
	}

	//pprof
	go http.ListenAndServe(":6868", nil)

	//monitor
	m := logic.NewMonitor()
	if err := m.Run(); err != nil {
		logrus.Fatalln(err)
	}

	//wait
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	s := <-ch
	log.Infof("receive signal[%s : %d], quit", s.String(), s)
}
