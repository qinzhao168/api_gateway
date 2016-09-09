package main

import (
	"api_gateway/basis/errors"
	"api_gateway/basis/etc"
	l "api_gateway/basis/log"
	"api_gateway/service"
	// "common/signal"
)

var log = l.New("api_gatway")

var HttpServerAddr = etc.String("api_gateway/uri", "server")

func main() {
	log.Info("build start: " + HttpServerAddr)
	go func() {
		if err := service.HttpServe(HttpServerAddr); err != nil {
			log.Fatal(errors.As(err))
		}
	}()

	// signal
	// signal.Serve()
	var signal = make(chan int, 1)
	<-signal
	log.Info("build end: " + HttpServerAddr)
}
