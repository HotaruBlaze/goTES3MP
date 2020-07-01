package main

import (
	"bytes"
	"encoding/json"
	"fmt"

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
	// sendMethod := s[0]
	res := &jsonResponce{}
	res.Method = s[0]
	// fmt.Println("You:", s)
	if viper.GetBool("debug") {
		log.Println("Length of array sent to relayProcess is", len(s))
	}
	switch res.Method {
	case "Tes3mp-Command":
		// onTes3mpCommand(s[3])
	case "IRC":
		log.Debugln("FROM IRC:", s[2])
		results := gjson.GetMany(s[2], "user", "method", "pid", "responce")

		PlayerName := results[0].String()
		Responce := results[3].String()

		// Send Message to discord
		DiscordSendMessage(PlayerName + ": " + Responce)

	case "Discord":
		res := &jsonResponce{}
		res.Method = "Discord"
		res.User = s[1]
		res.Pid = -1
		// res.Role = ""
		// res.RoleColor = ""
		res.Responce = s[3]
		fmt.Println("0:", s[0])
		fmt.Println("1:", s[1])
		fmt.Println("2:", s[2])
		fmt.Println("3:", s[3])
		// fmt.Println("4:", s[4])
		// fmt.Println("0:", s[0])
		log.Debugln("[Relay][Discord]", s)

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
		IRCSendMessage(sendResponce)
	default:
		log.Error(tes3mpLogMessage, `Something tried to use method "`+res.Method+`" but has no handler registered`)
	}
}
