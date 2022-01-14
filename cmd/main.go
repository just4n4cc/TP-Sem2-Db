package main

import (
	"github.com/just4n4cc/tp-sem2-db/internal/app"
	"github.com/just4n4cc/tp-sem2-db/pkg/logger"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"os"
)

const logMessage = "main:"

func main() {
	logger.Init(log.DebugLevel)

	//message := logMessage
	//defer logger.Info(message + "[EXIT]")
	application, err := app.NewApp()
	if err != nil {
		//logger.Error(message, err)
		os.Exit(1)
	}
	err = application.Run()
	if err != nil {
		//logger.Error(message, err)
		os.Exit(1)
	}
}
