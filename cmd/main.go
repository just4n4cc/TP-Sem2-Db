package main

import (
	"github.com/just4n4cc/tp-sem2-db/internal/app"
	log "github.com/just4n4cc/tp-sem2-db/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
)

const logMessage = "main:"

func main() {
	//log.Init(logrus.DebugLevel)
	log.Init(logrus.InfoLevel)

	message := logMessage
	defer log.Info(message + "[EXIT]")
	application, err := app.NewApp()
	if err != nil {
		log.Error(message, err)
		os.Exit(1)
	}
	err = application.Run()
	if err != nil {
		log.Error(message, err)
		os.Exit(1)
	}
}
