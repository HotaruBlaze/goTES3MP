package main

import (
	"bufio"
	"log"
	"os/exec"

	"github.com/spf13/viper"
)

// ServerRunning - TES3MP server Status
var ServerRunning string = "false"

// CurrentPlayers on TES3MP
var CurrentPlayers int = 0

// MaxPlayers placeholder variable
var MaxPlayers int = 0

// Version of goTES3MP
const Version = "v0.0.1"

var tes3mpLogMesage string = "[goTES3MP]: "

func main() {
	LoadConfig()
	go InitWebserver()
	go UpdateStatusTimer()
	go InitDiscord()
	LaunchTes3mp()

}

// LaunchTes3mp : Start and initialize TES3MP
func LaunchTes3mp() {
	tes3mpPath := viper.GetString("tes3mpPath")
	tes3mpBinary := "/tes3mp-server"

	cmd := exec.Command(tes3mpPath + tes3mpBinary)
	stdout, _ := cmd.StdoutPipe()

	startErr := cmd.Start()
	if startErr != nil {
		log.Fatalf("cmd.Run() failed with '%s'\n", startErr)
	}
	outScanner := bufio.NewScanner(stdout)
	for outScanner.Scan() {
		m := outScanner.Text()
		Linter(m)
	}
}
