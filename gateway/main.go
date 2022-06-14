package main

import (
	"gateway/infrastructure/api"
	"gateway/startup"
	"gateway/startup/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

func main() {
	api.Log = log
	log.Out = os.Stdout

	path := "gateway.log"
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(8760)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	if err == nil {
		log.SetOutput(writer)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	config := config.NewConfig()
	server, _ := startup.NewServer(config)
	server.Start()
	log.Info("Server staring...")
}
