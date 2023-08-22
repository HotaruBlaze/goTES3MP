package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var ipAddressArray []string

type IPHubResponseStruct struct {
	IP          string `json:"ip"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	Asn         int    `json:"asn"`
	Isp         string `json:"isp"`
	Block       int    `json:"block"`
}

type ipqualityscoreresponseStruct struct {
	Success        bool    `json:"success"`
	Message        string  `json:"message"`
	FraudScore     int     `json:"fraud_score"`
	CountryCode    string  `json:"country_code"`
	Region         string  `json:"region"`
	City           string  `json:"city"`
	ISP            string  `json:"ISP"`
	ASN            int     `json:"ASN"`
	Organization   string  `json:"organization"`
	IsCrawler      bool    `json:"is_crawler"`
	Timezone       string  `json:"timezone"`
	Mobile         bool    `json:"mobile"`
	Host           string  `json:"host"`
	Proxy          bool    `json:"proxy"`
	Vpn            bool    `json:"vpn"`
	Tor            bool    `json:"tor"`
	ActiveVpn      bool    `json:"active_vpn"`
	ActiveTor      bool    `json:"active_tor"`
	RecentAbuse    bool    `json:"recent_abuse"`
	BotStatus      bool    `json:"bot_status"`
	ConnectionType string  `json:"connection_type"`
	AbuseVelocity  string  `json:"abuse_velocity"`
	ZipCode        string  `json:"zip_code"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	RequestID      string  `json:"request_id"`
}

func checkPlayerIP(ipAddress string) bool {
	var wasIPBlocked bool

	if slices.Contains(ipAddressArray, ipAddress) {
		return true
	}

	// If no api keys are set, print out a warning and skip the checks.
	if len(viper.GetString("vpn.iphub_apikey")) == 0 && len(viper.GetString("vpn.ipqualityscore_apikey")) == 0 {
		log.Warnln("[vpnChecker]: ", "vpnChecker was triggered, however no api keys are currently set. Allowing player to join.")
		return false
	}

	// IPHub API Check
	if len(viper.GetString("vpn.iphub_apikey")) > 0 {
		wasIPBlocked = ipHubRequest(ipAddress)
		if wasIPBlocked {
			return true
		}
	}

	// IPQualityScore API Check
	if len(viper.GetString("vpn.ipqualityscore_apikey")) > 0 {
		wasIPBlocked = ipqualityscoreRequest(ipAddress)
	}

	return wasIPBlocked
}

func ipHubRequest(ipAddress string) bool {
	url := fmt.Sprintf("http://v2.api.iphub.info/ip/%s", ipAddress)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		checkError("ipHubRequest:1", err)
	}
	req.Header.Set("X-Key", viper.GetString("vpn.iphub_apikey"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		checkError("ipHubRequest:2", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		checkError("ipHubRequest:3", err)
	}

	var IPResponse IPHubResponseStruct
	err = json.Unmarshal(body, &IPResponse)
	if err != nil {
		checkError("ipHubRequest:4", err)
	}

	if IPResponse.Block == 1 {
		ipAddressArray = appendIfMissing(ipAddressArray, ipAddress)
		return true
	}
	return false
}

func ipqualityscoreRequest(ipAddress string) bool {
	apiKey := viper.GetString("vpn.ipqualityscore_apikey")
	webReq := fmt.Sprintf("https://ipqualityscore.com/api/json/ip/%s/%s", apiKey, ipAddress)

	req, err := http.NewRequest("GET", webReq, nil)
	if err != nil {
		checkError("ipqualityscoreRequest:1", err)
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		checkError("ipqualityscoreRequest:2", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		checkError("ipqualityscoreRequest:3", err)
		return false
	}

	var IPresponse ipqualityscoreresponseStruct
	err = json.Unmarshal(body, &IPresponse)
	if err != nil {
		checkError("ipqualityscoreRequest:4", err)
		return false
	}

	if IPresponse.FraudScore >= 80 {
		ipAddressArray = appendIfMissing(ipAddressArray, ipAddress)
		return true
	}

	return false
}
