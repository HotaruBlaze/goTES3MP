package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
)

// InitWebserver : Initalize webserver
func InitWebserver() {
	http.HandleFunc("/status", status)
	err := http.ListenAndServe(viper.GetString("webserver.port"), nil)
	if err != nil {
		log.Errorln("[Webserver]", "Unable to start webserver, %v", err)
	}
	time.Sleep(60 * time.Microsecond)
	log.Infoln("[Webserver] Started on port", viper.GetString("webserver.port"))

}

// status : Print current ServerStatus struct as json
func status(w http.ResponseWriter, r *http.Request) {
	s := UpdateStatus()
	if s == nil {
		log.Errorln("UpdateStatus returned nil")
		fmt.Fprint(w, "An Error Occurred while getting /status")
		return
	}
	status := pretty.Pretty(s)
	fmt.Fprint(w, status)
}
