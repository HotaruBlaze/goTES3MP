package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var ipAddressArray []string

type IPHubResponceStruct struct {
	IP          string `json:"ip"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	Asn         int    `json:"asn"`
	Isp         string `json:"isp"`
	Block       int    `json:"block"`
}

func checkPlayerIP(ipAddress string) int {
	if slices.Contains(ipAddressArray, ipAddress) {
		return 1
	}
	var webReq = "http://v2.api.iphub.info/ip/" + ipAddress
	req, err := http.NewRequest("GET", webReq, nil)
	if err != nil {
		checkError("checkPlayerIP:1", err)
	}
	req.Header.Set("X-Key", viper.GetString("iphub_apikey"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		checkError("checkPlayerIP:2", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		checkError("checkPlayerIP:3", err)
	}

	var IPResponce IPHubResponceStruct
	err = json.Unmarshal(body, &IPResponce)
	if err != nil {
		checkError("checkPlayerIP:4", err)
	}
	defer resp.Body.Close()

	if IPResponce.Block == 1 {
		ipAddressArray = AppendIfMissing(ipAddressArray, ipAddress)
		return 1
	}
	return 0
}
