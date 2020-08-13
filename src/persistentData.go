package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
		fmt.Println(err)
		os.Exit(1)
	}
	defer persistantDataFile.Close()

	byteValue, _ := ioutil.ReadAll(persistantDataFile)
	json.Unmarshal(byteValue, &persistantData)

}
func pdsaveData() {
	pd, err := json.Marshal(&persistantData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(pd))
	ioutil.WriteFile(persistantFilePath, pd, os.ModePerm)
}
