package main

import (
	"flag"
	"github.com/nate-martinez/faux-twitter/server/config"
	"github.com/nate-martinez/faux-twitter/server/log"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("c", "", "config path")

func main() {
	flag.Parse()

	if err := startup(); err != nil {
		os.Exit(1)
	}

	go listenAndServe()

	handleSigs()
}

func startup() error {
	if err := config.LoadFile(*configPath); err != nil {
		return err
	}
	log.Log.Debug("startup finished")
	return nil
}

func listenAndServe() {

}

func handleSigs() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	shutdown()
}

func shutdown() {}
