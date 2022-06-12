package main

import (
	"gateway/infrastructure/api"
	"gateway/startup"
	"gateway/startup/config"
	"github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func main() {
	api.Log = log
	log.Out = os.Stdout

	file, err := os.OpenFile("gateway.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	defer file.Close()

	config := config.NewConfig()
	server, _ := startup.NewServer(config)
	server.Start()
}
