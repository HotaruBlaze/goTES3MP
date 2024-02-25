package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type baseresponse struct {
	JobID    string            `json:"jobid"`
	ServerID string            `json:"serverid"`
	Method   string            `json:"method"`
	Source   string            `json:"source"`
	Target   string            `json:"target"`
	Data     map[string]string `json:"data"`
}

type commandResponse struct {
	JobID    string      `json:"jobid"`
	ServerID string      `json:"serverid"`
	Method   string      `json:"method"`
	Source   string      `json:"source"`
	Data     CommandData `json:"data"`
}

type CommandArg struct {
	Required    bool   `json:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CommandData struct {
	CommandArgs []*CommandArg `json:"args"`
	Command     string        `json:"command"`
	Description string        `json:"description"`
}

// handleIncomingMessage handles the incoming message and processes it accordingly.
//
// It takes a map[string]interface{} as the parameter and returns an interface{} and an error.
func handleIncomingMessage(data map[string]interface{}) (interface{}, error) {
	// Extract the method from the data
	method, ok := data["method"].(string)
	if !ok {
		return nil, errors.New("method is not a string")
	}

	// Process the message based on the method
	switch method {
	// Executing Discord slash command
	case "DiscordSlashCommand":
		// Print processing message if in debug mode
		if viper.GetBool("debug") {
			log.Println("Processing Discord Command:", method)
		}
		return data, nil
	// Registering Discord Slash Command
	case "RegisterDiscordSlashCommand":
		// Print registering message if in debug mode
		if viper.GetBool("debug") {
			log.Println("Registering a Discord Command:", method)
		}
		// Unmarshal the data into commandResponse struct
		var incomingData commandResponse
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonData, &incomingData)
		if err != nil {
			return nil, err
		}

		// Process the Discord command
		processDiscordCommand(incomingData)

		// Ensure commandResponses.Commands map is initialized
		if commandResponses.Commands == nil {
			commandResponses.Commands = make(map[string]CommandData)
		}

		// Add the new command to commandResponses
		newCommand := CommandData{
			Command:     incomingData.Data.Command,
			Description: incomingData.Data.Description,
			CommandArgs: incomingData.Data.CommandArgs,
		}
		commandResponses.Commands[incomingData.Data.Command] = newCommand

		// Create a string slice of argument names
		var argStrings []string
		for _, arg := range incomingData.Data.CommandArgs {
			argStrings = append(argStrings, arg.Name) // Or whatever field you want to pass as string
		}

		// Call AddDiscordCommand with the converted arguments
		AddDiscordCommand(&commandResponses, incomingData.Data.Command, newCommand.Description, argStrings...)

		return incomingData, nil
	default:
		// Print the method if it is not handled by specific logic above.
		// Decode the data into baseresponse struct and process the relay message
		var res baseresponse
		err := mapstructure.Decode(data, &res)
		if err != nil {
			return nil, err
		}
		processRelayMessage(res)
		return res, nil
	}
}

// processRelayMessage processes the relay message and returns a boolean.
//
// s: baseresponse object containing the relay message data
// returns: true if the relay message is processed successfully, false otherwise
func processRelayMessage(s baseresponse) bool {
	// Convert the baseresponse to a pointer
	res := &s

	// Perform sanity check on the relay message
	err := processRelayMessageSanityCheck(res)
	if err != nil {
		log.Errorln("processRelayMessageSanityCheck failed.")
		log.Errorf("Err %s", err)
		return false
	}

	// Log the length and data of the relay message if debug mode is enabled
	if viper.GetBool("debug") {
		log.Println("[Debug][processRelayMessage]:", len(res.Data))
		log.Println(res.Data)
	}

	// Check if ServerID is missing from the response
	if len(res.ServerID) == 0 {
		log.Warnln("ServerID Missing from response:")
		return false
	}

	// Process the relay message based on the method
	switch res.Method {
	case "Sync":
		serverSync(res.ServerID, res)
	case "IRC":
		log.Println("TODO: Method \"IRC\" Not Implemented Yet.")
	case "DiscordChat":
		// Send the relay message data to Discord chat
		jsonresponse, err := json.Marshal(res)
		checkError("DiscordChat", err)
		sendresponse := bytes.NewBuffer(jsonresponse).String()
		IRCSendMessage(viper.GetString("irc.systemchannel"), sendresponse)
		usrMsg := res.Data["User"] + ": " + res.Data["Message"]
		logRelayedMessages("Discord", usrMsg)
	case "rawDiscord":
		// Process the relay message for raw Discord
		var m rawDiscordStruct
		m.Channel = res.Data["channel"]
		m.Server = res.Data["server"]
		m.Message = res.Data["message"]

		// Check if required fields are not nil
		if m.Channel == "" || m.Server == "" || m.Message == "" {
			log.Errorln("[ProcessRelayMessage][rawDiscord]: One or more required fields are nil")
			return false
		} else {
			// Send the raw Discord message and log the relayed message
			status := sendRawDiscordMessage(m)
			logRelayedMessages("TES3MP", m.Message)
			return status
		}
	case "VPNCheck":
		// Process the VPN check for the relay message
		processVPNCheck(res)
	case "DiscordSlashCommandResponse":
		// Process the response for Discord slash command
		discordInteractiveToken := res.Data["discordInteractiveToken"]
		discordInteractiveReply := res.Data["response"]
		SendDiscordInteractiveMessage(discordInteractiveToken, discordInteractiveReply)
	default:
		log.Println(res.Method, " is an unknown method.")
	}

	return false
}

func processDiscordCommand(s commandResponse) bool {
	res := &s
	err := processRelayCommandSanityCheck(res)
	if err != nil {
		log.Errorln("processDiscordCommand failed.")
		log.Errorf("Err %s", err)
		return false
	}

	if len(res.ServerID) == 0 {
		log.Warnln("ServerID Missing from response:")
		return false
	}

	return true
}

// processVPNCheck is a function that processes VPN checks for a given response.
func processVPNCheck(res *baseresponse) {
	// Create a rawDiscordStruct from the data in the response
	m := rawDiscordStruct{
		Channel: res.Data["channel"],
		Server:  res.Data["server"],
		Message: res.Data["message"],
	}

	// Perform the VPN check on the player's IP
	isPlayerUsingVPN := checkPlayerIP(m.Message)

	// Set the kickPlayer field in the response data based on the result of the VPN check
	res.Data["kickPlayer"] = "no"
	if isPlayerUsingVPN {
		log.Printf("[VPNCheck]: %s has been kicked.", m.Message)
		res.Data["kickPlayer"] = "yes"
	} else {
		log.Printf("[VPNCheck]: %s is not suspected to be using a VPN.", m.Message)
	}

	// Convert the response to JSON
	jsonresponse, err := json.Marshal(res)
	checkError("VPNCheck", err)

	// Send the JSON response to the IRC system channel
	IRCSendMessage(viper.GetString("irc.systemchannel"), string(jsonresponse))
}

func logRelayedMessages(server string, message string) {
	if server != "" && message != "" {
		log.Println("<" + server + "> " + message)
	}
}

// processRelayMessageSanityCheck checks the sanity of the relay message.
// It ensures that the method is not blank and that data is provided.
// If any of the checks fail, it returns an error.
func processRelayMessageSanityCheck(relayMsg *baseresponse) error {
	// Check if the method is blank
	if relayMsg.Method == "" {
		return fmt.Errorf("method cannot be blank")
	}

	// Check if jobid is blank
	if relayMsg.JobID == "" {
		return fmt.Errorf("jobid cannot be blank")
	}

	// Check if data is provided
	if len(relayMsg.Data) == 0 {
		return fmt.Errorf("no data provided")
	}

	// Return nil if all checks pass
	return nil
}

// processRelayMessageSanityCheck checks the sanity of the relay message.
// It ensures that the method is not blank and that data is provided.
// If any of the checks fail, it returns an error.
func processRelayCommandSanityCheck(relayMsg *commandResponse) error {
	// Check if the method is blank
	if relayMsg.Method == "" {
		return fmt.Errorf("method cannot be blank")
	}

	// Check if jobid is blank
	if relayMsg.JobID == "" {
		return fmt.Errorf("jobid cannot be blank")
	}

	// Check if data is provided
	if relayMsg.Data.Command == "" {
		return fmt.Errorf("no command provided")
	}

	// Return nil if all checks pass
	return nil
}
