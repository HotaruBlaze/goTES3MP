package main

import (
	"bytes"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

type jsonResponce struct {
	Method    string `json:"method"`
	User      string `json:"user"`
	Pid       int    `json:"pid"`
	Role      string `json:"role"`
	RoleColor string `json:"role_color"`
	Responce  string `json:"responce"`
}

func relayProcess(s []string) {
	var systemChannel = viper.GetString("irc.systemchannel")
	var chatChannel = viper.GetString("irc.chatchannel")
	res := &jsonResponce{}
	res.Method = s[0]
	if viper.GetBool("debug") {
		log.Println("Length of array sent to relayProcess is", len(s))
		log.Println(s)
	}
	switch res.Method {
	case "Tes3mp-Command":
		// onTes3mpCommand(s[3])
	case "IRC":
		ircChannel := s[1]
		// Json System Message
		if ircChannel == systemChannel {
			results := gjson.GetMany(s[3], "user", "method", "pid", "responce")

			PlayerName := results[0].String()
			Responce := results[3].String()
			DiscordSendMessage(PlayerName + ": " + Responce)

			ircResponce := "[TES3MP]" + " " + PlayerName + ": " + Responce

			IRCSendMessage(chatChannel, ircResponce)
		}
		// From dedicated IRC Chat
		if ircChannel == chatChannel {
			res := &jsonResponce{}
			res.Method = "IRC"
			res.User = s[2]
			res.Pid = -1
			res.Responce = string(strings.Join(s[3:], " "))
			jsonResponce, err := json.Marshal(res)
			if err != nil {
				log.Errorln("[Relay]", "Failed to create JSON for chatChannel, ", err)
			}
			sendResponce := bytes.NewBuffer(jsonResponce).String()
			IRCSendMessage(systemChannel, sendResponce)
			DiscordSendMessage("[IRC] " + res.User + ": " + res.Responce)
		}

	case "Discord":
		res := &jsonResponce{}
		res.Method = "Discord"
		res.User = s[1]
		res.Pid = -1

		res.Responce = s[3]
		log.Debugln("[Relay][Discord]", s)

		// This does not work correctly if u have multiple valid roles
		// Not really sure how to fix this at this time.
		if len(s) > 4 {
			if s[4] != "" && s[5] != "" {
				res.Role = s[4]
				res.RoleColor = s[5]
			}
		}
		jsonResponce, err := json.Marshal(res)
		if err != nil {
			log.Error(tes3mpLogMessage, "Failed to generate jsonResponce for Discord->IRC.", "\n", "ERR: ", err)
		}

		sendResponce := bytes.NewBuffer(jsonResponce).String()
		IRCSendMessage(systemChannel, sendResponce)

		ircResponce := "[Discord] " + res.User + ": " + res.Responce
		IRCSendMessage(chatChannel, ircResponce)

	default:
		log.Error(tes3mpLogMessage, `Something tried to use method "`+res.Method+`" but has no handler registered`)
	}
}
