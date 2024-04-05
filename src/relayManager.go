package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	protocols "github.com/hotarublaze/gotes3mp/src/protocols"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// handleIncomingMessage handles the incoming message and processes it accordingly.
//
// It takes a map[string]interface{} as the parameter and returns an interface{} and an error.
func handleIncomingMessage(data *protocols.BaseResponse) (interface{}, error) {

	method := data.Method
	if len(method) == 0 {
		return nil, errors.New("method is not a string")
	}

	processRelayMessage(data)
	return "", nil
}

func handleIncomingComamnd(data protocols.DiscordSlashCommand) (interface{}, error) {
	method := data.Method
	if len(method) == 0 {
		return nil, errors.New("method is not a string")
	}

	// Process the message based on the method
	switch method {
	// Executing Discord slash command
	case "DiscordSlashCommand":
		// Convert this to the correct protocol
		var incomingData protocols.DiscordSlashCommand
		unmarshaler := jsonpb.Unmarshaler{}

		err := unmarshaler.Unmarshal(strings.NewReader(data.String()), &incomingData)
		if err != nil {
			return nil, err
		}

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
		// Process the Discord command
		VerifyDiscordCommand(data)

		// Ensure commandResponses.Commands map is initialized
		if commandResponses.Commands == nil {
			commandResponses.Commands = make(map[string]protocols.CommandData)
		}

		// Add the new command to commandResponses
		newCommand := protocols.CommandData{
			Command:     data.Data.Command,
			Description: data.Data.Description,
			Args:        data.Data.Args,
		}
		commandResponses.Commands[data.Data.Command] = newCommand

		// Create a string slice of argument names
		var argStrings []string
		for _, arg := range data.Data.Args {
			argStrings = append(argStrings, arg.Name) // Or whatever field you want to pass as string
		}

		// Call AddDiscordCommand with the converted arguments
		AddDiscordCommand(&commandResponses, data.Data.Command, data.Data.Description, argStrings...)

		return "", nil
	default:
		return "", nil
	}
}

// processRelayMessage processes the relay message and returns a boolean.
//
// s: baseresponse object containing the relay message data
// returns: true if the relay message is processed successfully, false otherwise
func processRelayMessage(s *protocols.BaseResponse) bool {
	// Convert the baseresponse to a pointer
	res := s

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
	if len(res.ServerId) == 0 {
		log.Warnln("ServerID Missing from response:")
		return false
	}

	// Process the relay message based on the method
	switch res.Method {
	case "Sync":
		serverSync(res.ServerId, res)
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
		var m protocols.RawDiscordStruct
		m.Channel = res.Data["channel"]
		m.Server = res.Data["server"]
		m.Message = res.Data["message"]

		// Check if required fields are not nil
		if m.Channel == "" || m.Server == "" || m.Message == "" {
			log.Errorln("[ProcessRelayMessage][rawDiscord]: One or more required fields are nil")
			return false
		} else {
			// Send the raw Discord message and log the relayed message
			status := sendRawDiscordMessage(&m)
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

func VerifyDiscordCommand(s protocols.DiscordSlashCommand) bool {

	err := processRelayCommandSanityCheck(&s)
	if err != nil {
		log.Errorln("VerifyDiscordCommand failed.")
		log.Errorf("Err %s", err)
		return false
	}

	if len(s.ServerId) == 0 {
		log.Warnln("ServerID Missing from response:")
		return false
	}

	return true
}

// processVPNCheck is a function that processes VPN checks for a given response.
func processVPNCheck(res *protocols.BaseResponse) {
	// Create a rawDiscordStruct from the data in the response
	m := protocols.RawDiscordStruct{
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
func processRelayMessageSanityCheck(relayMsg *protocols.BaseResponse) error {
	// Check if the method is blank
	if relayMsg.Method == "" {
		return fmt.Errorf("method cannot be blank")
	}

	// Check if jobid is blank
	if relayMsg.JobId == "" {
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
func processRelayCommandSanityCheck(relayMsg *protocols.DiscordSlashCommand) error {
	// Check if the method is blank
	if relayMsg.Method == "" {
		return fmt.Errorf("method cannot be blank")
	}

	// Check if jobid is blank
	if relayMsg.JobId == "" {
		return fmt.Errorf("jobid cannot be blank")
	}

	// Check if data is provided
	if relayMsg.Data.Command == "" {
		return fmt.Errorf("no command provided")
	}

	// Return nil if all checks pass
	return nil
}
