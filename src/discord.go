package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// InitDiscord Initialize discordgo
func InitDiscord() {
	Token := viper.GetString("discordToken")

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)
	discord.AddHandler(UpdateDiscordStatus)
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println(tes3mpLogMesage + "Discord Module is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

// UpdateDiscordStatus Update discord bot status to match player count on TES3MP
func UpdateDiscordStatus(s *discordgo.Session, event *discordgo.Ready) {
	for {
		var currentPlayers = strconv.Itoa(CurrentPlayers)
		var maxPlayers = strconv.Itoa(MaxPlayers)
		var status = ""

		serverName := viper.GetString("serverName")

		if len(serverName) > 0 {
			serverName = serverName + ": "
			status = serverName + currentPlayers + "/" + maxPlayers
			status = status + " Players"
		} else {
			status = currentPlayers + "/" + maxPlayers
			status = status + " Players"
		}

		idleSince := 0
		time.Sleep(5 * time.Second)
		usd := discordgo.UpdateStatusData{
			IdleSince: &idleSince,
			Game: &discordgo.Game{
				Name: status,
				Type: discordgo.GameTypeGame,
			},
			AFK:    false,
			Status: "online",
		}

		err := s.UpdateStatusComplex(usd)
		if err != nil {
			fmt.Println(err)
		}
	}
}
