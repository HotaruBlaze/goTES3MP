//go:generate protoc --go_out=paths=source_relative:../src --go_opt=paths=source_relative .\protocols\messages.proto
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GitCommit used to build
var GitCommit string

// Build version
var Build = "0.0.0-Debug"

// ServerRunning - TES3MP server Status
var ServerRunning bool

// CurrentPlayers on the server
var CurrentPlayers int

// MaxPlayers from .cfg file
var MaxPlayers int

// TES3MPVersion : Tes3mp Version
var TES3MPVersion = ""

// Players :  List is current Players Connected
var Players = []string{}

var tes3mpLogMessage = "[goTES3MP]"

// MultiWrite : Prints to logfile and os.Stdout
var MultiWrite io.Writer
var reader *bufio.Reader

func init() {
	if GitCommit == "" {
		GitCommit = "None"
	}
	initializeLogger()
	loadConfig()
	loadData()
}
func main() {
	debugEnabled := viper.GetBool("debug")
	if debugEnabled {
		log.Warnln("Debug mode is enabled")
		log.SetLevel(log.DebugLevel)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Infoln("Performing clean shutdown...")
		commandShutdown()
	}()

	if viper.GetBool("webserver.enable") {
		go InitWebserver()
	}
	go UpdateStatusTimer()
	if viper.GetBool("discord.enable") {
		go InitDiscord()
	}
	go InitIRC()
	if viper.GetBool("printMemoryInfo") {
		go MemoryDebugInfo()
	}

	if viper.GetBool("enableInteractiveConsole") {
		reader = bufio.NewReader(os.Stdin)
	}

	getStatus(true, false)
	for {
		time.Sleep(200 * time.Millisecond)

		if viper.GetBool("enableInteractiveConsole") {
			fmt.Print("> ")
			command, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			command = strings.TrimRight(command, "\r\n")
			args := strings.Split(command, " ")

			switch strings.ToLower(args[0]) {
			case "status":
				handleStatusCommand()
			case "reloadirc":
				handleReloadIRCCommand()
			case "reloaddiscord":
				log.Debugln("Attempting to reload Discord")
				InitDiscord()
			case "purgecommands":
				log.Println("Purging Discord commands...")
				purgeDiscordCommands()
			case "exit", "quit", "stop":
				log.Debugln("Shutting down...")
				commandShutdown()
			default:
				log.Warnf("Command " + command + " was not recognized.")
			}
		}
	}
}

func initializeLogger() {
	dt := time.Now()
	ProgramDirectory := "./goTES3MP/logs/"
	logfileName := ProgramDirectory + "goTES3MP-" + dt.Format("02-01-2006-15_04_05") + ".log"

	if _, err := os.Stat(ProgramDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(ProgramDirectory, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	logfile, err := os.OpenFile(logfileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	MultiWrite = io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(MultiWrite)

	log.SetLevel(log.InfoLevel)
}
