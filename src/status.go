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
	ServerOnline   string
	Version        string
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
	tes3mpPath := viper.GetString("tes3mpPath")

	props, err := ReadPropertiesFile(tes3mpPath + "/tes3mp-server-default.cfg")
	if err != nil {
		log.Println("Error while reading properties file")
	}
	if i, err := strconv.Atoi(props["maximumPlayers"]); err == nil {
		MaxPlayers = i
	}
}

// UpdateStatus for keeping server stats synced
func UpdateStatus() (s string) {
	CurrentStatus := ServerStatus{
		ServerOnline:   ServerRunning,
		Version:        Version,
		CurrentPlayers: CurrentPlayers,
		MaxPlayers:     MaxPlayers,
	}
	var jsonData []byte
	jsonData, err := json.Marshal(CurrentStatus)
	if err != nil {
		log.Println(err)
	}

	return string(jsonData)
}
