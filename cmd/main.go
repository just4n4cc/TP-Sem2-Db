package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	s := &http.Server{
		Addr: ":5000",
		Handler: ,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
