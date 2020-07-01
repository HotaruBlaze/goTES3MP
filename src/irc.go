package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	irc "github.com/thoj/go-ircevent"
)

var IrcServer string
var IrcPort string
var IrcNick string
var Channel string
var password string
var Irccon *irc.Connection

func InitIRC() {
	IrcServer = viper.GetString("irc.server")
	IrcPort = viper.GetString("irc.port")
	IrcNick = viper.GetString("irc.nick")
	Channel = viper.GetString("irc.channel")
	password = viper.GetString("irc.pass")
	// Irccon.Log = *Logger
	Irccon = irc.IRC(IrcNick, IrcNick)
	Irccon.Debug = false
	Irccon.UseTLS = false
	Irccon.AddCallback("001", func(e *irc.Event) {
		log.Infoln("[IRC] Connected to", IrcServer+":"+IrcServer, "as", IrcNick)
		Irccon.Join(Channel)
		log.Infoln("[IRC] Joined channel", Channel)
	})
	Irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			res := []string{"IRC", event.Nick, event.Message()}
			relayProcess(res)
		}(event)
	})
	err := Irccon.Connect(IrcServer + ":" + IrcPort)
	if err != nil {
		log.Printf("Err %s", err)
		return
	}
	Irccon.Join(Channel)
	log.Println(tes3mpLogMessage, "IRC Module is now running")

	Irccon.Loop()
}

// IRCSendMessage : Send message to IRC Channel
func IRCSendMessage(message string) {
	Irccon.Privmsg(Channel, message)
}
