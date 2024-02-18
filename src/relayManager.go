package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type baseresponse struct {
	JobID    string            `json:"jobid"`
	ServerID string            `json:"serverid"`
	Method   string            `json:"method"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	Data     map[string]string `json:"data"`
}

func processRelayMessage(s baseresponse) bool {
	res := &s
	err := processRelayMessageSanityCheck(res)
	if err != nil {
		log.Errorln("processRelayMessageSanityCheck failed.")
		log.Errorf("Err %s", err)
		return false
	}

	if viper.GetBool("debug") {
		log.Println("[Debug][processRelayMessage]:", len(res.Data))
		log.Println(res.Data)
	}

	if len(res.ServerID) == 0 {
		log.Warnln("ServerID Missing from response:")
		return false
	}

	switch res.Method {
	case "Sync":
		serverSync(res.ServerID, res)
	case "IRC":
		log.Println("TODO: Method \"IRC\" Not Implemented Yet.")
	case "DiscordChat":
		jsonresponse, err := json.Marshal(res)
		checkError("DiscordChat", err)
		sendresponse := bytes.NewBuffer(jsonresponse).String()
		IRCSendMessage(viper.GetString("irc.systemchannel"), sendresponse)
		usrMsg := res.Data["User"] + ": " + res.Data["Message"]
		logRelayedMessages("Discord", usrMsg)
	case "rawDiscord":
		var m rawDiscordStruct
		m.Channel = res.Data["channel"]
		m.Server = res.Data["server"]
		m.Message = res.Data["message"]

		// Check if channel,server and message are not nil, but print the reason why it failed.
		if m.Channel == "" || m.Server == "" || m.Message == "" {
			log.Errorln("[ProcessRelayMessage][rawDiscord]: One or more required fields are nil")
			return false
		} else {
			status := sendRawDiscordMessage(m)
			logRelayedMessages("TES3MP", m.Message)
			return status
		}
	case "VPNCheck":
		processVPNCheck(res)
	default:
		log.Println(res.Method, " is an unknown method.")
	}

	return false
}

// processVPNCheck is a function that processes VPN checks for a given response.
func processVPNCheck(res *baseresponse) {
	// Create a rawDiscordStruct from the data in the response
	m := rawDiscordStruct{
		Channel: res.Data["channel"],
		Server:  res.Data["server"],
		Message: res.Data["message"],
	}

	// Perform the VPN check on the player's IP
	isPlayerUsingVPN := checkPlayerIP(m.Message)

	// Set the kickPlayer field in the response data based on the result of the VPN check
	res.Data["kickPlayer"] = "no"
	if isPlayerUsingVPN {
		log.Printf("[VPNCheck]: %s has been kicked.", m.Message)
		res.Data["kickPlayer"] = "yes"
	} else {
		log.Printf("[VPNCheck]: %s is not suspected to be using a VPN.", m.Message)
	}

	// Convert the response to JSON
	jsonresponse, err := json.Marshal(res)
	checkError("VPNCheck", err)

	// Send the JSON response to the IRC system channel
	IRCSendMessage(viper.GetString("irc.systemchannel"), string(jsonresponse))
}

func logRelayedMessages(server string, message string) {
	if server != "" && message != "" {
		log.Println("<" + server + "> " + message)
	}
}

// processRelayMessageSanityCheck checks the sanity of the relay message.
// It ensures that the method is not blank and that data is provided.
// If any of the checks fail, it returns an error.
func processRelayMessageSanityCheck(relayMsg *baseresponse) error {
	// Check if the method is blank
	if relayMsg.Method == "" {
		return fmt.Errorf("method cannot be blank")
	}

	// Check if jobid is blank
	if relayMsg.JobID == "" {
		return fmt.Errorf("jobid cannot be blank")
	}

	// Check if data is provided
	if len(relayMsg.Data) == 0 {
		return fmt.Errorf("no data provided")
	}

	// Return nil if all checks pass
	return nil
}
