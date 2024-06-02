package main

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	protocols "github.com/hotarublaze/gotes3mp/src/protocols"
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
				// Parse a fuzzy metadata protocol to get the method
				var metadata protocols.Metadata
				metadata_unmarshaler := jsonpb.Unmarshaler{
					AllowUnknownFields: true,
				}
				unmarshaler := jsonpb.Unmarshaler{}

				err := metadata_unmarshaler.Unmarshal(strings.NewReader(event.Message()), &metadata)
				if err != nil {
					checkError("[IRC:AddCallback]: PRIVMSG 0", err)
				}

				switch metadata.Method {
				case "RegisterDiscordSlashCommand":
					{
						// Now parse this as a Discord Slash Command
						var dataPacket protocols.DiscordSlashCommand
						err := unmarshaler.Unmarshal(strings.NewReader(event.Message()), &dataPacket)
						if err != nil {
							checkError("[IRC:AddCallback]: PRIVMSG 1", err)
						} else {
							_, err := handleIncomingComamnd(dataPacket)
							if err != nil {
								checkError("[IRC:AddCallback]["+metadata.Method+"] PRIVMSG 2", err)
							}
						}
					}
				default:
					{
						// Now parse this as a normal system message
						var dataPacket protocols.BaseResponse
						err := unmarshaler.Unmarshal(strings.NewReader(event.Message()), &dataPacket)
						if err != nil {
							checkError("[IRC:AddCallback]: PRIVMSG 1", err)
						} else {
							_, err := handleIncomingMessage(&dataPacket)
							if err != nil {
								checkError("[IRC:AddCallback]: PRIVMSG 2", err)
							}
						}
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
		// This will hang if we forget this, but it's better than ignoring sig interrupt
		os.Exit(1)
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
