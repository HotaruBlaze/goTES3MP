package main

import (
	"bytes"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func discordCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	Data := make(map[string]string)
	var commandStruct baseResponce

	re := regexp.MustCompile(`[^\s"']+|([^\s"']*"([^"]*)"[^\s"']*)+|'([^']*)`)
	stringArr := re.FindAllString(m.Content[1:], -1)

	if viper.GetBool("debug") {
		log.Println("[Debug] discordCommandHandler:stringArr:'", stringArr)
	}

	commandStruct.ServerID = viper.GetViper().GetString("tes3mp.serverid")
	commandStruct.Method = "Command"
	commandStruct.Source = "DiscordCommand"
	Data["Command"] = stringArr[0]

	if len(stringArr) > 1 {
		Data["TargetPlayer"] = stringArr[1]
		if len(stringArr) > 2 {
			Data["CommandArgs"] = strings.Join(stringArr[2:], " ")
		}
	}

	Data["replyChannel"] = m.ChannelID
	commandStruct.Data = Data

	if viper.GetBool("debug") {
		log.Println("[Debug] discordCommandHandler:commandStruct'", Data)
	}

	log.Println("Staff Member '"+m.Author.Username+"' has executed the following command:", m.Content[1:])

	jsonResponce, err := json.Marshal(commandStruct)
	checkError("discordCommandHandler", err)
	sendResponce := bytes.NewBuffer(jsonResponce).String()
	IRCSendMessage(viper.GetString("irc.systemchannel"), sendResponce)
}
