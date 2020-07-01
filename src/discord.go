package main

import (
	"fmt"
	big "math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DiscordSession : Global Discord Session
var DiscordSession *discordgo.Session

// InitDiscord Initialize discordgo
func InitDiscord() {
	Token := viper.GetString("discord.token")

	Discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Errorln("error creating Discord session,", err)
		return
	}
	DiscordSession = Discord

	Discord.AddHandler(messageCreate)
	// Discord.StateEnabled = true
	Discord.AddHandler(UpdateDiscordStatus)
	err = Discord.Open()
	if err != nil {
		log.Errorln("error opening connection,", err)
		return
	}
	log.Println(tes3mpLogMessage, "Discord Module is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	Discord.Close()

}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	StaffRoles := viper.GetStringSlice("discord.staffRoles")
	if m.Author.ID == s.State.User.ID {
		return
	}
	// // If the message is "ping" reply with "Pong!"
	// if m.Content == "ping" {
	// 	s.ChannelMessageSend(m.ChannelID, "Pong!")
	// }

	// // If the message is "pong" reply with "Ping!"
	// if m.Content == "pong" {
	// 	s.ChannelMessageSend(m.ChannelID, "Ping!")
	// }
	if m.Content != "" && m.ChannelID == viper.GetString("discord.serverchat") {
		roleIndex := -1
		var HighestRole string
		var res []string
		var DiscordName string
		allowhexcolors := viper.GetBool("discord.allowcolorhexusage")
		fmt.Println("allowhexcolors = ", allowhexcolors)
		userroles := getDiscordRoles(m.Author.ID, m.GuildID)
		for role := range userroles {
			i, found := Find(StaffRoles, role)
			if found {
				if i == 0 || roleIndex < i {
					if viper.GetBool("debug") {
						log.Println(role, "ID:", i)
					}
					roleIndex = i
					HighestRole = role
				}
			}
		}
		gmember, _ := DiscordSession.GuildMember(m.GuildID, m.Author.ID)
		if gmember.Nick != "" {
			DiscordName = gmember.Nick
		} else {
			DiscordName = m.Author.Username
		}
		if roleIndex != -1 && HighestRole != "" {
			res = []string{"Discord", DiscordName, "-1", string(m.Content), HighestRole, userroles[HighestRole]}
			log.Debugln(res)
		} else {
			if allowhexcolors {
				res = []string{"Discord", DiscordName, "-1", string(m.Content), HighestRole, userroles[HighestRole]}
			} else {
				cleanString := stringVerifier(true, string(m.Content))
				res = []string{"Discord", DiscordName, "-1", string(cleanString)}

			}
		}
		// fmt.Println("Fuck:", res)

		relayProcess(res)
	}
}
func DiscordSendMessage(msg string) {
	serverChat := viper.GetString("discord.serverChat")

	DiscordSession.ChannelMessageSend(serverChat, msg)

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
			log.Errorln("UpdateDiscordStatus failed to update status.", err)
		}
	}
}

func getDiscordRoles(UserID string, GuildID string) map[string]string {
	roles := make(map[string]string)

	m, err := DiscordSession.GuildMember(GuildID, UserID)
	if err != nil {
		log.Errorln("ERR:", err)
	}
	for _, role := range m.Roles {
		role, err := DiscordSession.State.Role(GuildID, role)
		if err != nil {
			log.Errorln("getDiscordRoles Failed to get user roles from discord server.")
		}

		color := strings.ToUpper(
			toHexInt(
				big.NewInt(int64(role.Color)),
			),
		)
		roles[role.Name] = color
	}
	return roles
}
