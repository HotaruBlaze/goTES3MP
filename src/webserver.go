package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
)

// InitWebserver : Initalize webserver
func InitWebserver() {
	webport := viper.GetString("webserver.port")
	http.HandleFunc("/status", status)
	http.ListenAndServe(webport, nil)
}

// status : Print current ServerStatus struct as json
func status(w http.ResponseWriter, r *http.Request) {
	s := UpdateStatus()
	status := pretty.Pretty(s)
	if s == nil {
		log.Errorln("UpdateStatus returned nil")
		fmt.Fprintf(w, string("An Error Occurred while getting /status"))
	} else {
		fmt.Fprintf(w, string(status))
	}
}
