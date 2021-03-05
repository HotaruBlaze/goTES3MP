package main

import (
	"encoding/json"
	"io/ioutil"
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

// InitIRC : Initialize IRC
func InitIRC() {
	ircServer = viper.GetString("irc.server")
	ircPort = viper.GetString("irc.port")
	ircNick = viper.GetString("irc.nick")

	// IRC "System Channe;"
	systemchannel = viper.GetString("irc.systemchannel")
	// Add a extra channel for Talking via IRC
	chatchannel = viper.GetString("irc.chatchannel")

	password = viper.GetString("irc.pass")
	irccon = irc.IRC(ircNick, ircNick)
	irccon.Debug = false
	irccon.Log.SetOutput(ioutil.Discard)
	irccon.UseTLS = false
	irccon.Password = password
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
				var baseMsg baseResponce
				err := json.Unmarshal([]byte(event.Message()), &baseMsg)
				if err != nil {
					checkError("AddCallback: PRIVMSG", err)
				}
				processRelayMessage(baseMsg)
			}

		}(event)
	})
	err := irccon.Connect(ircServer + ":" + ircPort)
	if err != nil {
		log.Errorln("Failed to connect to IRC")
		log.Errorf("Err %s", err)
	}

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

	for {
		time.Sleep(10 * time.Second)
		if count < 6 {
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
		return
	}

}

// IRCSendMessage : Send message to IRC Channel
func IRCSendMessage(channel string, message string) {
	irccon.Privmsg(channel, message)
}
