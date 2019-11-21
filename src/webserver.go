package main

import (
	"fmt"
	"net/http"
)

// InitWebserver Start webfrontend
func InitWebserver() {
	http.HandleFunc("/debug", Debug)
	http.ListenAndServe(":8080", nil)
}

// Debug - Print current ServerStatus struct as json
func Debug(w http.ResponseWriter, r *http.Request) {
	s := UpdateStatus()
	fmt.Fprintf(w, s)
}
