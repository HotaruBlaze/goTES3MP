package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/spf13/viper"
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
		time.Sleep(5 * time.Second)
		_ = UpdateStatus()
	}

}
func getServerMaxPlayers() {
	tes3mpPath := viper.GetString("tes3mp.basedir")

	props, err := ReadPropertiesFile(tes3mpPath + "/tes3mp-server-default.cfg")
	if err != nil {
		log.Println("Error while reading properties file")
	}
	if i, err := strconv.Atoi(props["maximumPlayers"]); err == nil {
		MaxPlayers = i
	}
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
	// println(string(jsonData))
	// print(CurrentStatus)
	return jsonData
}
