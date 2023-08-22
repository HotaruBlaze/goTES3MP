package main

import (
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type persistantServerDataStruct struct {
	Users       map[string]string
	PlayerRoles []string `json:"PlayerRoles"`
}

var persistantData persistantServerDataStruct
var persistantFilePath = "./goTES3MP/data.json"

func loadData() {

	if _, err := os.Stat(persistantFilePath); os.IsNotExist(err) {
		saveData()
	}
	persistantDataFile, err := os.Open(persistantFilePath)
	if err != nil {
		log.Warnln("[loadData]: Command removerole errored with", err)
		os.Exit(1)
	}
	defer persistantDataFile.Close()

	byteValue, _ := io.ReadAll(persistantDataFile)
	err = json.Unmarshal(byteValue, &persistantData)
	if err != nil {
		log.Errorln("[loadData]", "Failed to Unmarshal Persistant Data, %v", err)
	}
}
func saveData() {
	pd, err := json.Marshal(&persistantData)
	if err != nil {
		log.Warnln("[pdsaveData]: Command removerole errored with", err)
		return
	}
	err = os.WriteFile(persistantFilePath, pd, os.ModePerm)
	if err != nil {
		log.Errorln("[pdsaveData]", "Failed to save file: , %v", err)
	}
}
