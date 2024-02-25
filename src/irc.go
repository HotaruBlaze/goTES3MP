package main

import (
	"encoding/json"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	irc "github.com/thoj/go-ircevent"
)

var ircServer string
var ircPort string
var ircNick string
var systemchannel string
var chatchannel string
var password string
var irccon *irc.Connection
var connectedToIRC bool

// InitIRC initializes the IRC connection using the configuration from viper
func InitIRC() {
	// Retrieve IRC configuration from viper
	ircServer = viper.GetString("irc.server")
	ircPort = viper.GetString("irc.port")
	ircNick = viper.GetString("irc.nick")

	// Define IRC channels and password
	systemchannel = viper.GetString("irc.systemchannel")
	chatchannel = viper.GetString("irc.chatchannel")
	password = viper.GetString("irc.pass")

	// Initialize IRC connection
	irccon = irc.IRC(ircNick, ircNick)
	irccon.Debug = false
	irccon.Log.SetOutput(io.Discard)
	irccon.UseTLS = false
	irccon.Password = password

	// Handle IRC connection events
	irccon.AddCallback("001", func(e *irc.Event) {
		log.Infoln("[IRC] Connected to", ircServer+":"+ircPort, "as", ircNick)
		irccon.Join(systemchannel)
		log.Infoln("[IRC] Joined channel", systemchannel)
		if chatchannel != "" && viper.GetBool("irc.enableChatChannel") {
			irccon.Join(chatchannel)
			log.Infoln("[IRC] Joined channel", chatchannel)
		}
	})
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			if event.Arguments[0] == systemchannel {
				var baseMsg map[string]interface{}
				err := json.Unmarshal([]byte(event.Message()), &baseMsg)
				if err != nil {
					checkError("[IRC:AddCallback]: PRIVMSG 1", err)
				} else {
					_, err := handleIncomingMessage(baseMsg)
					if err != nil {
						checkError("[IRC:AddCallback]: PRIVMSG 2", err)
					}
				}
			}
		}(event)
	})

	// Connect to IRC server
	err := irccon.Connect(ircServer + ":" + ircPort)
	if err != nil {
		log.Errorln("Failed to connect to IRC")
		log.Errorf("Err %s", err)
	}

	// Start IRC loop
	log.Println(tes3mpLogMessage, "IRC Module is now running")
	connectedToIRC = true
	irccon.Loop()
}

func ircReconnect() {
	count := 0
	irccon.QuitMessage = "Bot Restarting"

	currentstatus := irccon.Connected()
	if currentstatus {
		log.Println(tes3mpLogMessage, "[IRC] Shutting down Module")
		irccon.Quit()
		connectedToIRC = false
		log.Println(tes3mpLogMessage, "[IRC] Module Offline")
		time.Sleep(5 * time.Second)
	}
	log.Println(tes3mpLogMessage, "[IRC] Module Loading...")

	for count < 6 {
		time.Sleep(10 * time.Second)
		currentstatus := irccon.Connected()
		if !currentstatus {
			connectedToIRC = false
			count++
			err := irccon.Reconnect()
			if err != nil {
				log.Fatal(err)
			} else {
				connectedToIRC = true
			}
		}

		if connectedToIRC && irccon.Connected() {
			log.Println(tes3mpLogMessage, "[IRC] Module online...")
			return
		}
	}
	log.Error("Unable to Reconnect to IRC within 60 seconds.")
}

// IRCSendMessage : Send message to IRC Channel
func IRCSendMessage(channel string, message string) {
	irccon.Privmsg(channel, message)
}
