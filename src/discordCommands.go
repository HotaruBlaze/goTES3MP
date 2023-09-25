package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// discordCommandHandler handles the command received from Discord.
// It takes a session and a message create event as parameters.
func discordCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Split the message content into an array of strings
	stringArr := strings.Fields(m.Content[1:])

	// Create a command struct with the necessary fields
	commandStruct := baseresponse{
		ServerID: viper.GetViper().GetString("tes3mp.serverid"),
		Method:   "Command",
		Source:   "DiscordCommand",
		Data: map[string]string{
			"Command": stringArr[0],
		},
	}

	// Check if there are additional arguments in the command
	if len(stringArr) > 1 {
		commandStruct.Data["CommandArgs"] = strings.Join(stringArr[1:], "^")
	}

	// Set the reply channel in the command struct
	commandStruct.Data["replyChannel"] = m.ChannelID

	// Print debug information if debug mode is enabled
	if viper.GetBool("debug") {
		log.Println("[Debug] discordCommandHandler:commandStruct'", commandStruct.Data)
	}

	// Log the executed command
	log.Println("Staff Member '"+m.Author.Username+"' has executed the following command:", m.Content[1:])

	// Convert the command struct to JSON
	jsonresponse, err := json.Marshal(commandStruct)
	checkError("discordCommandHandler", err)
	sendresponse := string(jsonresponse)

	// Send the JSON response to the IRC system channel
	IRCSendMessage(viper.GetString("irc.systemchannel"), sendresponse)
}
