package main

import (
	"os"

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
	irccon.UseTLS = false
	irccon.Password = password
	irccon.AddCallback("001", func(e *irc.Event) {
		log.Infoln("[IRC] Connected to", ircServer+":"+ircPort, "as", ircNick)
		irccon.Join(systemchannel)
		log.Infoln("[IRC] Joined channel", systemchannel)
		irccon.Join(chatchannel)
		log.Infoln("[IRC] Joined channel", chatchannel)
	})
	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			if event.Arguments[0] == systemchannel {
				res := []string{"IRC", event.Arguments[0], event.Nick, event.Message()}
				relayProcess(res)
			}
			if event.Arguments[0] == chatchannel {
				res := []string{"IRC", event.Arguments[0], event.Nick, event.Message()}
				relayProcess(res)
			}
		}(event)
	})
	err := irccon.Connect(ircServer + ":" + ircPort)
	if err != nil {
		log.Errorln("Failed to connect to IRC")
		log.Errorf("Err %s", err)
		os.Exit(1)
	}
	log.Println(tes3mpLogMessage, "IRC Module is now running")

	irccon.Loop()
}

// IRCSendMessage : Send message to IRC Channel
func IRCSendMessage(channel string, message string) {
	irccon.Privmsg(channel, message)
}
