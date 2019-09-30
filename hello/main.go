package main

import (
	"bytes"
	"encoding/json"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-plugins/server/http"
	"github.com/segdumping/micro/registry"
	"github.com/segdumping/shared/signal"
	http2 "net/http"
	"os"
)

func Hello(writer http2.ResponseWriter, request *http2.Request) {
	log.Log("http endpoint hello")
	request.ParseForm()
	name := request.Form.Get("name")
	if len(name) == 0 {
		name = "world"
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode("hello" + " " + name)
	writer.Write(buf.Bytes())
}

func Greet(writer http2.ResponseWriter, request *http2.Request) {
	log.Log("http endpoint greet")
	request.ParseForm()
	name := request.Form.Get("name")
	if len(name) == 0 {
		name = "anonymous person"
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode("greet" + " " + name)
	writer.Write(buf.Bytes())
}

func main() {
	//create server with options
	srv := http.NewServer(
		server.Id("1"),
		server.Name("go.micro.api.hello"),
		server.Registry(registry.NewRegistry()),
		)

	//route
	mux := http2.NewServeMux()
	mux.HandleFunc("/hello/hello", Hello)
	mux.HandleFunc("/hello/greet", Greet)

	hd := srv.NewHandler(mux,
		server.EndpointMetadata("/hello/hello", map[string]string{"handler":"http"}),
		server.EndpointMetadata("/hello/greet", map[string]string{"handler":"http"}),
	)

	//register handler
	if err := srv.Handle(hd); err != nil {
		log.Fatal(err)
	}


	if err := srv.Init(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Error(err)
		}
	}()

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch)
	s := <-ch
	log.Infof("receive signal[%s : %d], quit", s.String(), s)
}
