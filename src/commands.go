package main

import (
	"os"

	color "github.com/fatih/color"
	"github.com/spf13/viper"
)

func commandStatus() {
	getStatus(false, true)
}
func commandShutdown() {
	if viper.GetBool("debug") {
		color.HiBlack("[DEBUG][shutdown] Done.")
	}
	irccon.QuitMessage = "Bot shutting down"
	irccon.Quit()
	DiscordSession.Close()
	os.Exit(0)
}

func commandIrcReconnect() {
	ircReconnect()
}
