package main

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// PlayerStruct : the Struct used for json SYSTEM messages
type PlayerStruct struct {
	Name     string `json:"user"`
	Method   string `json:"method"`
	Responce string `json:"responce"`
	Pid      int    `json:"pid"`
}

func tes3mpOutputHandler(s string) {
	ServerOutput := viper.GetBool("serveroutput")
	if ServerOutput {
		fmt.Println(s)
	}

	var isTrue bool
	isTrue, _ = regexp.MatchString(`TES3MP dedicated server`, s)
	if isTrue {
		dedicatedServerInfo := strings.Split(s, " ")
		build := dedicatedServerInfo[4] + " " + dedicatedServerInfo[5]
		build = strings.Replace(strings.Replace(build, ")", "", -1), "(", "", -1)
		if dedicatedServerInfo[3] == "------------------------------------------------------------" {
			log.Warnln(tes3mpLogMessage, "Malformed Version was caught")
			TES3MPVersion = "Unknown"
		} else {
			log.Infoln(tes3mpLogMessage, "TES3MP Version:", dedicatedServerInfo[3], build)
			TES3MPVersion = dedicatedServerInfo[3] + " " + build
		}

	}
	isTrue, _ = regexp.MatchString(`Called "OnServerPostInit"`, s)
	if isTrue {
		log.Infoln(tes3mpLogMessage, "Tes3mp server is now online")
		ServerRunning = true
		DiscordSendAlert("**" + "[TES3MP] Server is online. :yellow_heart:" + "**")
	}
	isTrue, _ = regexp.MatchString(`Called "OnServerExit"`, s)
	if isTrue {
		log.Infoln(tes3mpLogMessage, "Tes3mp server is now Offline")
		ServerRunning = false
		DiscordSendAlert("**" + "[TES3MP] Server is offline. :warning:" + "**")
	}
}
