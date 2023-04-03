package main

import (
	"bytes"
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type baseResponce struct {
	ServerID string            `json:"serverid"`
	Method   string            `json:"method"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	Data     map[string]string `json:"data"`
}

func processRelayMessage(s baseResponce) bool {
	var isValid bool
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

	if len(res.ServerID) > 0 {
		if viper.GetBool("debug") {
			log.Println("[Debug]:", "ServerID found:", res.ServerID)
		}
		isValid = true
	} else {
		log.Warnln("ServerID Missing from Responce:")
	}
	if isValid {
		switch res.Method {
		case "Sync":
			serverSync(res.ServerID, res)
			return false
		case "IRC":
			log.Println("TODO: Method \"IRC\" Not Implemented Yet.")
			return false
		case "DiscordChat":
			jsonResponce, err := json.Marshal(res)
			checkError("DiscordChat", err)
			sendResponce := bytes.NewBuffer(jsonResponce).String()
			IRCSendMessage(viper.GetString("irc.systemchannel"), sendResponce)
			usrMsg := res.Data["User"] + ": " + res.Data["Message"]
			logRelayedMessages("Discord", usrMsg)
		case "Discord":
			log.Println("Replaced with rawDiscord, Depreciation of this method is planned")
			return false
		case "rawDiscord":
			var m rawDiscordStruct
			m.Channel = res.Data["channel"]
			m.Server = res.Data["server"]
			m.Message = res.Data["message"]
			status := rawDiscordMessage(m)
			logRelayedMessages("TES3MP", m.Message)
			return status
		case "VPNCheck":
			var m rawDiscordStruct
			m.Channel = res.Data["channel"]
			m.Server = res.Data["server"]
			m.Message = res.Data["message"]
			blockLevel := checkPlayerIP(m.Message)
			if blockLevel == 1 {
				log.Println("[VPNCheck]:", m.Message, "has been kicked.")
				res.Data["kickPlayer"] = "yes"
			} else {
				log.Println("[VPNCheck]:", m.Message, "is not suspected to be using a VPN.")
				res.Data["kickPlayer"] = "no"
			}
			jsonResponce, err := json.Marshal(res)
			checkError("VPNCheck", err)
			sendResponce := bytes.NewBuffer(jsonResponce).String()
			IRCSendMessage(viper.GetString("irc.systemchannel"), sendResponce)
		default:
			log.Println(res.Method, " is an unknown method.")
		}
	}
	return false
}

func logRelayedMessages(server string, message string) {
	if server != "" && message != "" {
		log.Println("<" + server + "> " + message)
	}
}

func processRelayMessageSanityCheck(Rmsg *baseResponce) error {
	tempRelayMsg := Rmsg
	// Tried to convert this to a switch, it didnt like it.
	if tempRelayMsg.Method == "" {
		return errors.New("processRelayMessage: method cannot be blank")
	}
	if len(tempRelayMsg.Data) == 0 {
		return errors.New("processRelayMessage: No data provided.")
	}
	return nil
}
