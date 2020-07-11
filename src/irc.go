package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	irc "github.com/thoj/go-ircevent"
)

var IrcServer string
var IrcPort string
var IrcNick string
var Systemchannel string
var Chatchannel string
var password string
var Irccon *irc.Connection

func InitIRC() {
	IrcServer = viper.GetString("irc.server")
	IrcPort = viper.GetString("irc.port")
	IrcNick = viper.GetString("irc.nick")
	// Talking back and forth via json
	Systemchannel = viper.GetString("irc.systemchannel")
	// Add a extra channel for Talking via IRC
	Chatchannel = viper.GetString("irc.chatchannel")
	password = viper.GetString("irc.pass")
	Irccon = irc.IRC(IrcNick, IrcNick)
	Irccon.Debug = false
	Irccon.UseTLS = false
	Irccon.AddCallback("001", func(e *irc.Event) {
		log.Infoln("[IRC] Connected to", IrcServer+":"+IrcPort, "as", IrcNick)
		Irccon.Join(Systemchannel)
		log.Infoln("[IRC] Joined channel", Systemchannel)
		Irccon.Join(Chatchannel)
		log.Infoln("[IRC] Joined channel", Chatchannel)
	})
	Irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			if event.Arguments[0] == Systemchannel {
				res := []string{"IRC", event.Arguments[0], event.Nick, event.Message()}
				relayProcess(res)
				fmt.Println("This is a System Message")
			}
			if event.Arguments[0] == Chatchannel {
				res := []string{"IRC", event.Arguments[0], event.Nick, event.Message()}
				relayProcess(res)
				fmt.Println("This is a Chat Message")

			}
		}(event)
	})
	err := Irccon.Connect(IrcServer + ":" + IrcPort)
	if err != nil {
		log.Errorln("Failed to connect to IRC")
		log.Errorf("Err %s", err)
		os.Exit(1)
		// return
	}
	// Irccon.Join(Channel)
	log.Println(tes3mpLogMessage, "IRC Module is now running")

	Irccon.Loop()
}

// IRCSendMessage : Send message to IRC Channel
func IRCSendMessage(channel string, message string) {
	Irccon.Privmsg(channel, message)
}
