package main

import (
	"fmt"
	"github.com/spf13/viper"
	"regexp"
)

//Linter Lint output for Tes3mp Server
func Linter(s string) bool {
	var foundrule = false
	tes3mpServerUp, _ := regexp.MatchString(`Called "OnServerPostInit"`, s)
	tes3mpPlayerLeave, _ := regexp.MatchString(`Called "OnPlayerDisconnect"`, s)
	tes3mpPlayerJoin, _ := regexp.MatchString(`Called "OnPlayerConnect"`, s)

	if tes3mpServerUp && foundrule == false {
		fmt.Println(tes3mpLogMesage + "Tes3mp server is now online")
		ServerRunning = "true"
		foundrule = true
	}
	if tes3mpPlayerLeave && foundrule == false {
		isInvalidPlayer, _ := regexp.MatchString(`Unlogged player `, s)
		if isInvalidPlayer {

		} else {
			fmt.Println(tes3mpLogMesage + "Player Disconnected")
			CurrentPlayers = CurrentPlayers - 1
		}
		foundrule = true
	}
	if tes3mpPlayerJoin && foundrule == false {
		fmt.Println(tes3mpLogMesage + "Player Joined")
		CurrentPlayers = CurrentPlayers + 1
		foundrule = true
	}
	debugToggle := viper.GetBool("debug")

	if debugToggle {
		fmt.Println(s)
	}
	return foundrule
}
