package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func discordCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	Data := make(map[string]string)
	var commandStruct baseResponce
	commandStruct.ServerID = viper.GetViper().GetString("tes3mp.serverid")
	commandStruct.Method = "Command"
	commandStruct.Source = "DiscordCommand"
	Data["Command"] = m.Content[1:]
	Data["replyChannel"] = m.ChannelID
	commandStruct.Data = Data

	log.Println("Staff Member '"+m.Author.Username+"' has executed the following command:", m.Content[1:])

	jsonResponce, err := json.Marshal(commandStruct)
	checkError("discordCommandHandler", err)
	sendResponce := bytes.NewBuffer(jsonResponce).String()
	IRCSendMessage(viper.GetString("irc.systemchannel"), sendResponce)
}
