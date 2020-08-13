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

type discordRole struct {
	Position int    `json:"position"`
	Color    string `json:"color"`
	Name     string `json:"name"`
}

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
	defer InitDiscord()
	DiscordSession = Discord

	Discord.AddHandler(messageCreate)
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
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID == viper.GetString("discord.serverchat") {
		if viper.GetBool("discord.enablecommands") {
			if strings.HasPrefix(m.Content, viper.GetString("discord.commandprefix")) {
				messageCommand(m)
				return
			}
		}
	}
	if m.Content != "" && m.ChannelID == viper.GetString("discord.serverchat") {
		allowhexcolors := viper.GetBool("discord.allowcolorhexusage")
		var staffMember bool = false
		var userroles []string
		var res []string
		discordRoles := getDiscordRoles(m.Author.ID, m.GuildID)
		index, pos := -1, -1
		for r, i := range discordRoles {
			userroles = append(userroles, i.Name)
			_, validRole := FindinArray(persistantData.PlayerRoles, i.Name)
			if validRole {
				if i.Name == persistantData.Users[m.Author.ID] {
					index = r
					pos = i.Position
					break
				} else {
					if i.Position > pos {
						index = r
						pos = i.Position
					}
				}
			}
		}

		var DiscordName string
		gmember, _ := DiscordSession.GuildMember(m.GuildID, m.Author.ID)
		if gmember.Nick != "" {
			DiscordName = gmember.Nick
		} else {
			DiscordName = m.Author.Username
		}

		if index != -1 && discordRoles[index].Name != "" {
			res = []string{"Discord", DiscordName, "-1", string(m.Content), discordRoles[index].Name, discordRoles[index].Color}
			log.Debugln(res)
		} else {
			if allowhexcolors || staffMember {
				res = []string{"Discord", DiscordName, "-1", string(m.Content), discordRoles[index].Name, discordRoles[index].Color}
			} else {
				cleanString := stringVerifier(true, string(m.Content))
				res = []string{"Discord", DiscordName, "-1", string(cleanString)}
			}
		}
		relayProcess(res)
	}
}

func messageCommand(m *discordgo.MessageCreate) {
	commandArgs := strings.Split(trimLeftChar(m.Content), " ")
	commandArgs[0] = strings.ToLower(commandArgs[0])
	fmt.Println(commandArgs)
	if commandArgs[0] == "addrole" {
		if isStaffMember(m.Author.ID, m.GuildID) {
			if len(commandArgs) > 1 && commandArgs[1] != "" {
				_, i := FindinArray(persistantData.PlayerRoles, commandArgs[1])
				if i == false {
					g, err := DiscordSession.Guild(m.GuildID)
					if err != nil {
						fmt.Println(err)
						return
					}
					r := g.Roles
					var roles []string
					for _, v := range r {
						roles = AppendIfMissing(roles, v.Name)
					}
					_, roleExists := FindinArray(roles, commandArgs[1])
					if roleExists {
						persistantData.PlayerRoles = append(persistantData.PlayerRoles, commandArgs[1])
						pdsaveData()
						DiscordSendMessage("`Added role " + commandArgs[1] + " as avalable role`")
					} else {
						DiscordSendMessage("`Role " + commandArgs[1] + " does not exist on discord." + "`")

					}
				} else {
					DiscordSendMessage("`Role " + commandArgs[1] + " already Exists`")
				}
			} else {
				DiscordSendMessage("`Role name Missing.`")
			}
		} else {
			DiscordSendMessage("`You do not have permission for this command.`")
		}
		return
	}
	if commandArgs[0] == "removerole" {
		if isStaffMember(m.Author.ID, m.GuildID) {
			if len(commandArgs) > 1 && commandArgs[1] != "" {
				_, i := FindinArray(persistantData.PlayerRoles, commandArgs[1])
				if i == true {
					g, err := DiscordSession.Guild(m.GuildID)
					if err != nil {
						fmt.Println(err)
						return
					}
					r := g.Roles
					var roles []string
					for _, v := range r {
						roles = AppendIfMissing(roles, v.Name)
					}
					_, roleExists := FindinArray(roles, commandArgs[1])
					if roleExists {
						persistantData.PlayerRoles = RemoveEntryFromArray(persistantData.PlayerRoles, commandArgs[1])
						for n, r := range persistantData.Users {
							if r == commandArgs[1] {
								delete(persistantData.Users, n)
							}
						}
						pdsaveData()
						DiscordSendMessage("`Removed role " + commandArgs[1] + "`")
					} else {
						DiscordSendMessage("`Role " + commandArgs[1] + " does not exist on discord." + "`")

					}
				} else {
					DiscordSendMessage("`Role " + commandArgs[1] + " does not Exists`")
				}
			} else {
				DiscordSendMessage("`Role name Missing.`")
			}
		} else {
			DiscordSendMessage("`You do not have permission for this command.`")
		}
		return
	}
	if commandArgs[0] == "setrole" {
		var chatRoles, userRoles, validRoles []string
		// Get all avalable PlayerRoles
		for _, role := range persistantData.PlayerRoles {
			chatRoles = AppendIfMissing(chatRoles, role)
		}

		// Get users roles
		serverRoles := getDiscordRoles(m.Author.ID, m.GuildID)
		for _, r := range serverRoles {
			userRoles = AppendIfMissing(userRoles, r.Name)
		}
		// Get roles User is allowed to use
		for _, v := range chatRoles {
			_, found := FindinArray(userRoles, v)
			if found {
				validRoles = AppendIfMissing(validRoles, v)
			}
		}
		if len(commandArgs) > 1 {
			_, f := FindinArray(validRoles, commandArgs[1])
			if f {
				if persistantData.Users == nil {
					persistantData.Users = map[string]string{
						m.Author.ID: commandArgs[1],
					}
				}
				persistantData.Users[m.Author.ID] = commandArgs[1]
				pdsaveData()
				DiscordSendMessage("`Set role to " + commandArgs[1] + "`")
			} else {
				fmt.Println("You do not have access to this role.")
			}

		} else {
			var msg string
			msg = msg + "You have the following roles avalable." + "\n" + "```" + "\n"
			for _, r := range validRoles {
				msg = msg + r + "\n"
			}
			msg = msg + "```"
			if persistantData.Users[m.Author.ID] != "" {
				msg = msg + "You current role is "
				msg = msg + "`" + persistantData.Users[m.Author.ID] + "`"
			}
			fmt.Println(msg)
			DiscordSendMessage(msg)
		}
	}
	if commandArgs[0] == "list" {
		responce := ""
		if len(Players) > 0 {
			responce = "```" + "\n"
			for _, player := range Players {
				responce = responce + player + "\n"
			}
			responce = responce + "\n" + "```"
		}
		responce = responce + "`" + "Currently online: " + strconv.Itoa(len(Players)) + "/" + strconv.Itoa(MaxPlayers) + "`"
		DiscordSendMessage(responce)
	}
}

// DiscordSendMessage : Send message to Discord serverChat
func DiscordSendMessage(msg string) {
	serverChat := viper.GetString("discord.serverChat")
	DiscordSession.ChannelMessageSend(serverChat, msg)
}

// DiscordSendAlert : Send Alert to Discord alertsChannel
func DiscordSendAlert(msg string) {
	alertChannel := viper.GetString("discord.alertsChannel")
	DiscordSession.ChannelMessageSend(alertChannel, msg)
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
			log.Warnln("UpdateDiscordStatus failed to update status.", err)
		}
		time.Sleep(5 * time.Second)
	}
}

func getDiscordRoles(UserID string, GuildID string) []discordRole {
	var discordRoles []discordRole
	member, err := DiscordSession.GuildMember(GuildID, UserID)
	if err != nil {
		log.Errorln("ERR:", err)
	}
	for _, role := range member.Roles {
		role, err := DiscordSession.State.Role(GuildID, role)
		if err != nil {
			log.Errorln("getDiscordRoles Failed to get user roles from discord server.")
		}
		color := strings.ToUpper(
			toHexInt(
				big.NewInt(int64(role.Color)),
			),
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
		println("Failed to figure out if a user is a staff member.")
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
