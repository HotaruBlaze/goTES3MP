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

// InitDiscord initializes the discordgo session
func InitDiscord() {
	// Create a new discordgo session
	discord, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		log.Errorln("error creating Discord session:", err)
		return
	}
	defer discord.Close()

	// Set the global Discord session variable
	DiscordSession = discord

	// Add event handlers
	discord.AddHandler(messageCreate)
	discord.AddHandler(ready)
	discord.AddHandler(UpdateDiscordStatus)

	// Open the connection to Discord
	if err := discord.Open(); err != nil {
		log.Errorln("error opening connection:", err)
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

// allowColorHexUsage checks if the usage of color hex codes is allowed.
// It returns true if hex code usage is allowed or if the user is a staff member.
func allowColorHexUsage(msg *discordgo.Message) bool {
	// Check if hex code usage is allowed
	allowHexColors := viper.GetBool("discord.allowcolorhexusage")
	if allowHexColors {
		return true
	}

	// Check if the user is a staff member
	isStaff := isStaffMember(msg.Author.ID, msg.GuildID)
	return isStaff
}

// UpdateDiscordStatus updates the status of a Discord bot
func UpdateDiscordStatus(
	s *discordgo.Session,
	event *discordgo.Ready,
) {
	// Create a ticker that ticks every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Convert currentPlayers and maxPlayers to strings
		currentPlayers := strconv.Itoa(CurrentPlayers)
		maxPlayers := strconv.Itoa(MaxPlayers)

		// Initialize the status string
		status := ""

		// Get the serverName from the configuration file
		serverName := viper.GetString("serverName")

		if len(serverName) > 0 {
			// Add the serverName to the status if it is not empty
			serverName = serverName + ": "
			status = serverName + currentPlayers + "/" + maxPlayers
			status = status + " Players"
		} else {
			// Otherwise, only include the currentPlayers and maxPlayers
			status = currentPlayers + "/" + maxPlayers
			status = status + " Players"
		}

		idleSince := 0

		// Create the UpdateStatusData struct
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
			// Update the status using the Discord session
			err := s.UpdateStatusComplex(usd)
			if err != nil {
				// If there is an error, handle it accordingly
				if err == discordgo.ErrWSNotFound {
					log.Println("no websocket connection exists, Attempting reconnection in 5 seconds")
					_ = DiscordSession.Open()
				}
				log.Warnln("UpdateDiscordStatus failed to update status.", err)
			}
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

// isStaffMember checks if a user is a staff member based on their roles in a Discord guild.
func isStaffMember(UserID string, GuildID string) bool {
	// Get the list of staff roles from the configuration
	staffRoles := viper.GetStringSlice("discord.staffroles")

	// Get the roles of the user in the specified guild
	discordRoles := getDiscordRoles(UserID, GuildID)

	// Get the guild information
	guild, err := DiscordSession.Guild(GuildID)
	if err != nil {
		// Log a warning if there was an error fetching the guild information
		log.Warnln("isStaffMember", "Failed to figure out if a user is a staff member.")
	}

	// Check if the user is the owner of the guild
	if UserID == guild.OwnerID {
		return true
	}

	// Check if any of the user's roles match the staff roles
	for _, i := range discordRoles {
		_, found := FindinArray(staffRoles, i.Name)
		if found {
			return true
		}
	}

	// Return false if the user is not a staff member
	return false
}
