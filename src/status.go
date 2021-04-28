package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	color "github.com/fatih/color"
)

// ServerStatus struct
type ServerStatus struct {
	ServerOnline   bool
	TES3MPVersion  string
	Build          string
	Players        []string
	CurrentPlayers int
	MaxPlayers     int
}

// UpdateStatusTimer run UpdateStatus loop
func UpdateStatusTimer() {
	getServerMaxPlayers()
	for {
		time.Sleep(20 * time.Second)
		_ = UpdateStatus()
	}

}

func getServerMaxPlayers() {
	MaxPlayers = 0
}

// UpdateStatus for keeping server stats synced
func UpdateStatus() (s []byte) {
	CurrentStatus := ServerStatus{
		ServerOnline:   ServerRunning,
		TES3MPVersion:  TES3MPVersion,
		Build:          Build,
		Players:        Players,
		CurrentPlayers: CurrentPlayers,
		MaxPlayers:     MaxPlayers,
	}
	var jsonData []byte
	jsonData, err := json.Marshal(CurrentStatus)
	if err != nil {
		log.Println(err)
	}
	return jsonData
}
func getStatus(firstLaunch bool, showModules bool) {
	if firstLaunch {
		color.HiBlack(strings.Repeat("=", 32))
	}
	color.HiBlack("goTES3MP: " + Build)
	color.HiBlack("Commit: " + GitCommit)
	color.HiBlack("Github: " + "https://github.com/hotarublaze/goTES3MP" + "\n")
	if firstLaunch {
		color.HiBlack(strings.Repeat("=", 32))
	}
}
