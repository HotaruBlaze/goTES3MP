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

func rawDiscordMessage(rawDiscordStruct rawDiscordStruct) bool {
	_, err := DiscordSession.ChannelMessageSend(rawDiscordStruct.Channel, rawDiscordStruct.Message)
	checkError("rawDiscordMessage", err)
	return true
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if len(m.Content) > 0 {
		if m.Content[:1] == viper.GetString("discord.commandprefix") && isStaffMember(m.Author.ID, m.GuildID) {
			discordCommandHandler(s, m)
			return
		}
	}
	if m.ChannelID != viper.GetString("discord.serverchat") {
		return
	}
	var discordResponce baseResponce

	discordResponce.ServerID = viper.GetViper().GetString("tes3mp.serverid")
	discordResponce.Method = "DiscordChat"
	discordResponce.Source = "Discord"
	discordResponce.Target = "TES3MP"
	var user, message string

	if !allowcolorhexusage(m.Message) {
		message = removeRGBHex(m.Content)
	} else {
		message = m.Content
	}
	guildMember, err := s.GuildMember(m.GuildID, m.Message.Author.ID)
	checkError("[RelayDiscord]: guildMember ", err)
	hasNickname := guildMember.Nick

	if len(hasNickname) > 0 {
		user = hasNickname
	} else {
		user = m.Author.Username
	}

	var discordData map[string]string
	roleName, roleColor := getUsersRole(m.Message)
	if len(roleName) > 0 && len(roleColor) > 0 {
		discordData = map[string]string{
			"User":      user,
			"Message":   message,
			"RoleName":  roleName,
			"RoleColor": roleColor, // last comma is a must
		}
	} else {
		discordData = map[string]string{
			"User":      user,
			"Message":   message,
			"RoleName":  "",
			"RoleColor": "", // last comma is a must
		}

	}

	discordResponce.Data = discordData
	processRelayMessage(discordResponce)
}

func getUsersRole(m *discordgo.Message) (string, string) {
	discordRoles := getDiscordRoles(m.Author.ID, m.GuildID)
	var userroles []string
	var allowedUserRole, allowedStaffRole bool
	index, pos := -1, -1
	for r, i := range discordRoles {
		userroles = append(userroles, i.Name)
		if len(viper.GetViper().GetStringSlice("discord.userroles")) > 0 {
			discordUserroles := viper.GetViper().GetStringSlice("discord.userroles")
			_, allowedUserRole = FindinArray(discordUserroles, i.Name)
		}

		if len(viper.GetViper().GetStringSlice("discord.staffroles")) > 0 {
			discordStaffroles := viper.GetViper().GetStringSlice("discord.staffroles")
			_, allowedStaffRole = FindinArray(discordStaffroles, i.Name)
		}
		if i.Name == persistantData.Users[m.Author.ID] {
			index = r
			pos = i.Position
			break
		} else {
			if i.Position > pos && allowedStaffRole {
				index = r
				pos = i.Position
			} else {
				if i.Position > pos && allowedUserRole {
					index = r
					pos = i.Position
				}
			}
		}
	}
	if index == -1 {
		return "", ""
	} else {
		return discordRoles[index].Name, discordRoles[index].Color
	}
}
