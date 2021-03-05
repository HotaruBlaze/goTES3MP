package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type persistantServerDataStruct struct {
	Users       map[string]string
	PlayerRoles []string `json:"PlayerRoles"`
}

var persistantData persistantServerDataStruct
var persistantFilePath = "./goTES3MP/data.json"

func pdloadData() {

	if _, err := os.Stat(persistantFilePath); os.IsNotExist(err) {
		pdsaveData()
	}
	persistantDataFile, err := os.Open(persistantFilePath)
	if err != nil {
		log.Warnln("[pdloadData]: Command removerole errored with", err)
		os.Exit(1)
	}
	defer persistantDataFile.Close()

	byteValue, _ := ioutil.ReadAll(persistantDataFile)
	err = json.Unmarshal(byteValue, &persistantData)
	if err != nil {
		log.Errorln("[pdloadData]", "Failed to Unmarshal Persistant Data, %v", err)
	}
}
func pdsaveData() {
	pd, err := json.Marshal(&persistantData)
	if err != nil {
		log.Warnln("[pdsaveData]: Command removerole errored with", err)
		return
	}
	err = ioutil.WriteFile(persistantFilePath, pd, os.ModePerm)
	if err != nil {
		log.Errorln("[pdsaveData]", "Failed to save file: , %v", err)
	}
}
