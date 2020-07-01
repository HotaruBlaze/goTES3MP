package main

// WIP

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
)

// InitWebserver Start webfrontend
func InitWebserver() {
	http.HandleFunc("/status", status)
	http.ListenAndServe(":8080", nil)
}

// serverInfo - Print current ServerStatus struct as json
func status(w http.ResponseWriter, r *http.Request) {
	s := UpdateStatus()
	status := pretty.Pretty(s)
	if s == nil {
		log.Errorln("Oh god, what did u do")
	}
	fmt.Fprintf(w, string(status))
}
