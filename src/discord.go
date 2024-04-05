package main

import (
	"bytes"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	protocols "github.com/hotarublaze/gotes3mp/src/protocols"
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
var DiscordGuildID string
var defaultMemberPermissions int64 = discordgo.PermissionManageServer
var DMPermission bool = false

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
	discord.AddHandler(handleDiscordCommands)

	// Open the connection to Discord
	if err := discord.Open(); err != nil {
		log.Errorln("error opening connection:", err)
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Check if bot is in any discord servers first
	if len(event.Guilds) == 0 {
		log.Errorln("[Discord] Bot is not in any Discord servers")
		os.Exit(1)
	}
	if len(event.Guilds) > 1 {
		log.Warnln("[Discord] Bot is in more than 1 Discord server, this can have unintended results.")
	}

	// Set the playing status.
	err := s.UpdateGameStatus(0, "")
	if err != nil {
		log.Println(err)
	} else {
		// Discord module is ready!
		log.Println(tes3mpLogMessage, "Discord Module is now running")
		// Get the first guildID
		// DiscordGuildID = event.Guilds[0].ID
		DiscordGuildID = viper.GetString("discord.guildID")
		// Load Commands
		commandResponses, err = LoadDiscordCommandData()
		if err != nil {
			log.Errorln("Error loading Discord command data:", err)
		}
	}
}

func handleDiscordCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the interaction is an application command
	if i.Type == discordgo.InteractionApplicationCommand {
		// Get the name of the command
		commandName := strings.ToLower(i.ApplicationCommandData().Name)

		// Convert discord options to json we can handle easier
		commandArgs, err := discordOptionsToJSON(i.ApplicationCommandData().Options)
		if err != nil {
			log.Errorln("Error converting Discord options to JSON:", err)
		}

		commandArgs = string(commandArgs)

		// Find and execute the corresponding functionality based on the command name
		_, ok := commandResponses.Commands[commandName]
		if ok {
			// Build a DiscordCommand packet for TES3MP
			discordCommand := &protocols.BaseResponse{
				JobId:    uuid.New().String(),
				ServerId: viper.GetString("tes3mp.serverid"),
				Method:   "Command",
				Source:   "DiscordCommand",
				Target:   "TES3MP",
				Data: map[string]string{
					"command":                 commandName,
					"commandArgs":             commandArgs,
					"discordInteractiveToken": string(i.Interaction.Token),
				},
			}
			jsonresponse, err := json.Marshal(discordCommand)
			checkError("DiscordChat", err)
			sendresponse := bytes.NewBuffer(jsonresponse).String()
			IRCSendMessage(viper.GetString("irc.systemchannel"), sendresponse)

			// Temp response for now
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Processing...",
				},
			})
		} else {
			// Respond with unknown command message
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Unknown command. Type `/help` to see available commands.",
				},
			})
		}
	}
}

// discordOptionsToJSON converts Discord options to JSON format.
func discordOptionsToJSON(options []*discordgo.ApplicationCommandInteractionDataOption) (string, error) {
	// Create a map to store the data
	data := make(map[string]interface{})

	// Iterate through each option and convert it to the appropriate type
	for _, option := range options {
		name := option.Name
		var value interface{}
		switch option.Type {
		case discordgo.ApplicationCommandOptionString:
			value = option.StringValue()
		case discordgo.ApplicationCommandOptionInteger:
			value = option.IntValue()
		case discordgo.ApplicationCommandOptionBoolean:
			value = option.BoolValue()
		default:
			value = option.StringValue()
		}

		// Add the converted value to the data map
		data[name] = value
	}

	// Marshal the data map into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// createSlashCommand creates a new slash command for the given command name in the Discord guild.
// It maps each command argument to a discordgo.ApplicationCommandOption and sets the options for the command.
func createSlashCommand(command string) error {
	// Retrieve the command details from the commandResponses map
	tes3mpCommand := commandResponses.Commands[command]

	// Create a slice to hold the options
	var options []*discordgo.ApplicationCommandOption

	// Map each command argument to a discordgo.ApplicationCommandOption
	for _, arg := range tes3mpCommand.Args {
		// Determine the type of the argument based on your requirements
		optionType := discordgo.ApplicationCommandOptionString // For example, assuming all args are strings

		// Create the option
		option := &discordgo.ApplicationCommandOption{
			Type:        optionType,
			Name:        arg.Name,
			Description: arg.Description,
			Required:    arg.Required,
		}

		// Add the option to the slice
		options = append(options, option)
	}

	// Define the data for the slash command
	commandData := &discordgo.ApplicationCommand{
		Name:                     tes3mpCommand.Command,
		Description:              tes3mpCommand.Description,
		Type:                     discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: &defaultMemberPermissions,
		DMPermission:             &DMPermission,
		Options:                  options, // Set the options for the command
	}

	// Create the slash command in a specific guild
	_, err := DiscordSession.ApplicationCommandCreate(DiscordSession.State.User.ID, DiscordGuildID, commandData)
	if err != nil {
		return err
	}

	// Print confirmation message
	log.Println("Created discord slash command:", command)
	return nil
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

func SendDiscordInteractiveMessage(interactionToken, newContent string) {
	// Construct the interaction response data for editing
	responseEdit := &discordgo.WebhookEdit{
		Content: &newContent, // New content for the interaction response
	}

	// Update the interaction response
	DiscordSession.WebhookMessageEdit(DiscordSession.State.User.ID, interactionToken, "@original", responseEdit)
}
