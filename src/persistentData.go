package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type persistantDataStruct struct {
	Users       map[string]string
	PlayerRoles []string `json:"PlayerRoles"`
}

var persistantData persistantDataStruct
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
	json.Unmarshal(byteValue, &persistantData)

}
func pdsaveData() {
	pd, err := json.Marshal(&persistantData)
	if err != nil {
		log.Warnln("[pdsaveData]: Command removerole errored with", err)
		return
	}
	ioutil.WriteFile(persistantFilePath, pd, os.ModePerm)
}
