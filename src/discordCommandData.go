package main

import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type CommandResponses struct {
	Commands map[string]CommandData `json:"commands"`
}

var commandResponses CommandResponses

// AddDiscordCommand is responsible for adding a new Discord slash command to the system. It saves the command and registers it with the Discord platform.
func AddDiscordCommand(data *CommandResponses, command string, description string, args ...string) {
	// If the data has no Commands map, return without doing anything
	if data.Commands == nil {
		return
	}

	// Create CommandArg objects for each argument
	commandArgs := make([]*CommandArg, len(args))
	for i, arg := range args {
		commandArgs[i] = &CommandArg{
			Required:    true,
			Name:        arg,
			Description: description,
		}
	}

	// Add the new CommandData to the Commands map in the data
	data.Commands[command] = CommandData{
		Command:     command,
		Description: description,
		CommandArgs: commandArgs,
	}

	// Create a new slash command
	createSlashCommand(command)

	// Save the updated Discord command data to a file
	if err := SaveDiscordCommandData(*data, "discordCommands.json"); err != nil {
		// Print an error message if there was an error saving the data
		log.Errorln("Error saving Discord command data:", err)
	}
}

// RemoveDiscordCommand removes the specified command from the CommandResponses map.
//
// data *CommandResponses - the map of command responses
// command string - the command to be removed
func RemoveDiscordCommand(data *CommandResponses, command string) {
	delete(data.Commands, command)
}

// SaveDiscordCommandData saves the given CommandResponses data to a file specified by the filename parameter.
// It takes a CommandResponses data and a filename string as parameters and returns an error.
func SaveDiscordCommandData(data CommandResponses, filename string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %v", err)
	}
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON to file: %v", err)
	}
	return nil
}

// LoadDiscordCommandData loads the Discord command data from a JSON file.
//
// It returns a CommandResponses and an error.
func LoadDiscordCommandData() (CommandResponses, error) {
	fileData, err := os.ReadFile("discordCommands.json")
	if err != nil {
		return CommandResponses{}, fmt.Errorf("error reading file: %v", err)
	}
	var loadedData CommandResponses
	err = json.Unmarshal(fileData, &loadedData)
	if err != nil {
		return CommandResponses{}, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	// Count the number of registered commands
	numCommands := len(loadedData.Commands)
	log.Println("[Discord]", "Loaded", numCommands, "Discord commands!")

	return loadedData, nil
}

func purgeDiscordCommands() error {
	commands, err := DiscordSession.ApplicationCommands(DiscordSession.State.User.ID, DiscordGuildID)
	if err != nil {
		return err
	}

	for _, command := range commands {
		err := DiscordSession.ApplicationCommandDelete(DiscordSession.State.User.ID, DiscordGuildID, command.ID)
		if err != nil {
			return err
		}
		RemoveDiscordCommand(&commandResponses, command.Name)
		log.Println("[Discord]", "Purged Discord command:", command.Name)
	}

	log.Println("[Discord]", "Purged all Discord commands!")
	SaveDiscordCommandData(commandResponses, "discordCommands.json")
	return nil
}
