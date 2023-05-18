package main

import (
	"math/big"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type discordRole struct {
	Position int    `json:"position"`
	Color    string `json:"color"`
	Name     string `json:"name"`
}

// DiscordSession : Global Discord Session
var DiscordSession *discordgo.Session

// InitDiscord Initialize discordgo
func InitDiscord() {
	Discord, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.Errorln("error creating Discord session,", err)
		return
	}
	defer Discord.Close()

	DiscordSession = Discord
	Discord.AddHandler(messageCreate)
	Discord.AddHandler(ready)
	Discord.AddHandler(UpdateDiscordStatus)

	err = Discord.Open()
	if err != nil {
		log.Errorln("error opening connection,", err)
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	err := s.UpdateGameStatus(0, "")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(tes3mpLogMessage, "Discord Module is now running")
	}
}

func allowcolorhexusage(message *discordgo.Message) bool {
	allowhexcolors := viper.GetBool("discord.allowcolorhexusage")
	if allowhexcolors {
		return true
	}

	isStaff := isStaffMember(message.Author.ID, message.GuildID)
	return isStaff
}

// UpdateDiscordStatus: Update discord bot status
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
		usd := discordgo.UpdateStatusData{
			IdleSince: &idleSince,
			Activities: []*discordgo.Activity{{
				Name: status,
				Type: discordgo.ActivityTypeGame,
			}},
			AFK:    false,
			Status: "online",
		}
		if s.DataReady {
			err := s.UpdateStatusComplex(usd)
			if err != nil {
				if err == discordgo.ErrWSNotFound {
					log.Println("no websocket connection exists, Attempting reconnection in 5 seconds")
					time.Sleep(5 * time.Second)
					_ = DiscordSession.Open()
				}
				log.Warnln("UpdateDiscordStatus failed to update status.", err)
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func getDiscordRoles(UserID string, GuildID string) []discordRole {
	var discordRoles []discordRole
	member, err := DiscordSession.GuildMember(GuildID, UserID)
	if err != nil {
		checkError("getDiscordRoles", err)
	}
	for _, role := range member.Roles {
		role, err := DiscordSession.State.Role(GuildID, role)
		if err != nil {
			log.Errorln("getDiscordRoles Failed to get user roles from discord server.")
		}
		color := toHexInt(
			big.NewInt(int64(role.Color)),
		)
		r := discordRole{
			Position: role.Position,
			Color:    color,
			Name:     role.Name,
		}
		discordRoles = append(discordRoles, r)
	}
	return discordRoles
}

func isStaffMember(UserID string, GuildID string) bool {
	staffRoles := viper.GetStringSlice("discord.staffroles")
	discordRoles := getDiscordRoles(UserID, GuildID)
	guild, err := DiscordSession.Guild(GuildID)
	if err != nil {
		log.Warnln("isStaffMember", "Failed to figure out if a user is a staff member.")
	}
	if UserID == guild.OwnerID {
		return true
	}
	for _, i := range discordRoles {
		_, found := FindinArray(staffRoles, i.Name)
		if found {
			return true
		}
	}
	return false
}
