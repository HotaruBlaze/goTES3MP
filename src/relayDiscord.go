package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type rawDiscordStruct struct {
	Channel string `json:"channel"`
	Server  string `json:"server"`
	Message string `json:"Message"`
}

func sendRawDiscordMessage(rawDiscordStruct rawDiscordStruct) bool {
	_, err := DiscordSession.ChannelMessageSend(rawDiscordStruct.Channel, rawDiscordStruct.Message)
	checkError("rawDiscordMessage", err)
	return true
}

// messageCreate is a function that handles incoming messages in a Discord server
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || len(m.Content) == 0 || m.ChannelID != viper.GetString("discord.serverchat") {
		return
	}

	if m.Content[:1] == viper.GetString("discord.commandprefix") && isStaffMember(m.Author.ID, m.GuildID) {
		discordCommandHandler(s, m)
		return
	}

	discordresponse := baseresponse{
		ServerID: viper.GetViper().GetString("tes3mp.serverid"),
		Method:   "DiscordChat",
		Source:   "Discord",
		Target:   "TES3MP",
	}

	message := m.Content
	if !allowColorHexUsage(m.Message) {
		message = removeRGBHex(message)
	}
	message = convertDiscordFormattedItems(message, m.GuildID)
	message = filterDiscordEmotes(message)

	guildMember, err := s.GuildMember(m.GuildID, m.Message.Author.ID)
	checkError("[RelayDiscord]: guildMember ", err)
	user := guildMember.Nick
	if user == "" {
		user = m.Author.Username
	}

	roleName, roleColor := getUsersRole(m.Message)

	discordData := map[string]string{
		"User":      user,
		"Message":   message,
		"RoleName":  roleName,
		"RoleColor": roleColor,
	}

	discordresponse.Data = discordData

	processRelayMessage(discordresponse)
}

func getUsersRole(m *discordgo.Message) (string, string) {
	discordRoles := getDiscordRoles(m.Author.ID, m.GuildID)
	var allowedUserRole, allowedStaffRole bool
	index := -1

	for r, i := range discordRoles {
		discordUserroles := viper.GetViper().GetStringSlice("discord.userroles")
		if len(discordUserroles) > 0 {
			_, allowedUserRole = FindinArray(discordUserroles, i.Name)
		}

		discordStaffroles := viper.GetViper().GetStringSlice("discord.staffroles")
		if len(discordStaffroles) > 0 {
			_, allowedStaffRole = FindinArray(discordStaffroles, i.Name)
		}

		if i.Name == persistantData.Users[m.Author.ID] {
			index = r
			break
		} else if i.Position > discordRoles[index].Position && allowedStaffRole {
			index = r
		} else if i.Position > discordRoles[index].Position && allowedUserRole {
			index = r
		}
	}

	if index == -1 {
		return "", ""
	}

	return discordRoles[index].Name, discordRoles[index].Color
}
