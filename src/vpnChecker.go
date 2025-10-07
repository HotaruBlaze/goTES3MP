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
	VPN            bool    `json:"vpn"`
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

	log.Debugf("[vpnChecker] Starting IP check for: %s", ipAddress)

	if slices.Contains(ipAddressArray, ipAddress) {
		log.Infof("[vpnChecker] IP %s is in blocked list, denying access", ipAddress)
		return true
	}

	// If no api keys are set, print out a warning and skip the checks.
	if len(viper.GetString("vpn.iphub_apikey")) == 0 && len(viper.GetString("vpn.ipqualityscore_apikey")) == 0 {
		log.Warnln("[vpnChecker]: ", "vpnChecker was triggered, however no api keys are currently set. Allowing player to join.")
		return false
	}

	// IPHub API Check
	if len(viper.GetString("vpn.iphub_apikey")) > 0 {
		log.Debugf("[vpnChecker] Checking IP %s with IPHub API", ipAddress)
		wasIPBlocked = ipHubRequest(ipAddress)
		if wasIPBlocked {
			log.Infof("[vpnChecker] IP %s blocked by IPHub API", ipAddress)
			return true
		}
		log.Debugf("[vpnChecker] IP %s not blocked by IPHub API", ipAddress)
	}

	// IPQualityScore API Check
	if len(viper.GetString("vpn.ipqualityscore_apikey")) > 0 {
		log.Debugf("[vpnChecker] Checking IP %s with IPQualityScore API", ipAddress)
		wasIPBlocked = ipqualityscoreRequest(ipAddress)
		if wasIPBlocked {
			log.Infof("[vpnChecker] IP %s blocked by IPQualityScore API", ipAddress)
		} else {
			log.Debugf("[vpnChecker] IP %s not blocked by IPQualityScore API", ipAddress)
		}
	}

	log.Debugf("[vpnChecker] Final result for IP %s: blocked=%v", ipAddress, wasIPBlocked)
	return wasIPBlocked
}

func ipHubRequest(ipAddress string) bool {
	url := fmt.Sprintf("http://v2.api.iphub.info/ip/%s", ipAddress)
	log.Debugf("[vpnChecker] IPHub request URL: %s", url)

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

	log.Debugf("[vpnChecker] IPHub response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		checkError("ipHubRequest:3", err)
	}

	log.Debugf("[vpnChecker] IPHub response body: %s", string(body))

	var IPResponse IPHubResponseStruct
	err = json.Unmarshal(body, &IPResponse)
	if err != nil {
		checkError("ipHubRequest:4", err)
	}

	log.Debugf("[vpnChecker] IPHub parsed response - Block: %d, Country: %s, ISP: %s",
		IPResponse.Block, IPResponse.CountryName, IPResponse.Isp)

	if IPResponse.Block == 1 {
		log.Infof("[vpnChecker] IP %s blocked by IPHub (Block=1)", ipAddress)
		ipAddressArray = appendIfMissing(ipAddressArray, ipAddress)
		return true
	}

	log.Debugf("[vpnChecker] IP %s not blocked by IPHub (Block!=1)", ipAddress)
	return false
}

func ipqualityscoreRequest(ipAddress string) bool {
	apiKey := viper.GetString("vpn.ipqualityscore_apikey")
	webReq := fmt.Sprintf("https://ipqualityscore.com/api/json/ip/%s/%s", apiKey, ipAddress)
	log.Debugf("[vpnChecker] IPQualityScore request URL: %s", webReq)

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

	log.Debugf("[vpnChecker] IPQualityScore response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		checkError("ipqualityscoreRequest:3", err)
		return false
	}

	log.Debugf("[vpnChecker] IPQualityScore response body: %s", string(body))

	var ipqs ipqualityscoreresponseStruct
	err = json.Unmarshal(body, &ipqs)
	if err != nil {
		checkError("ipqualityscoreRequest:4", err)
		return false
	}

	log.Debugf("[vpnChecker] IPQualityScore parsed response - VPN: %v, Tor: %v, Proxy: %v, FraudScore: %d, Country: %s",
		ipqs.VPN, ipqs.Tor, ipqs.Proxy, ipqs.FraudScore, ipqs.CountryCode)

	// ---- Decision logic ----
	if ipqs.VPN {
		log.Infof("[vpnChecker] IP %s blocked by IPQualityScore (VPN=true)", ipAddress)
		return true
	}

	if ipqs.Tor {
		log.Infof("[vpnChecker] IP %s blocked by IPQualityScore (Tor=true)", ipAddress)
		return true
	}

	// Proxy + very high fraud score
	if ipqs.Proxy && ipqs.FraudScore >= 95 {
		log.Infof("[vpnChecker] IP %s blocked by IPQualityScore (Proxy=true, FraudScore=%d)", ipAddress, ipqs.FraudScore)
		return true
	}

	// Medium risk -> warn only
	if ipqs.Proxy && ipqs.FraudScore >= 85 {
		log.Warnf("[VpnChecker] Suspicious IP (proxy=true, fraud=%d): %s", ipqs.FraudScore, ipAddress)
		return false
	}

	log.Debugf("[vpnChecker] IP %s not blocked by IPQualityScore", ipAddress)
	return false
}
