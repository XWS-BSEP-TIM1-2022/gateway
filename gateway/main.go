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
	/*Log rotation correlation function
	`Withlinkname 'establishes a soft connection for the latest logs
	`Withrotationtime 'sets the time of log splitting, and how often to split
	Only one of withmaxage and withrotationcount can be set
	 `Withmaxage 'sets the maximum save time before cleaning the file
	 `Withrotationcount 'sets the maximum number of files to be saved before cleaning
	*/
	//The following configuration logs rotate a new file every 1 minute, keep the log files of the last 3 minutes, and automatically clean up the surplus.
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
