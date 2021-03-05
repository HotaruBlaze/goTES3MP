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

	color "github.com/fatih/color"
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

func init() {
	if len(GitCommit) == 0 {
		GitCommit = "None"
	}
	initLogger()
	LoadConfig()
	pdloadData()
}
func main() {
	enableDebug := viper.GetBool("debug")
	if enableDebug {
		log.Warnln("Debug mode is enabled")
		log.SetLevel(log.DebugLevel)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Infoln("Preforming clean shutdown...")
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

	reader := bufio.NewReader(os.Stdin)

	getStatus(true, false)
	for {
		time.Sleep(2 * 100 * time.Millisecond)
		// TODO: This should be tweaked so ">" is always at the bottom.
		fmt.Print("> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		command = strings.TrimSuffix(command, "\n")

		args := strings.Split(command, " ")
		switch strings.ToLower(args[0]) {
		case "status":
			commandStatus()
		case "reloadirc":
			commandIrcReconnect()
		case "exit", "quit", "stop":
			color.HiBlack("Shutting down...")
			commandShutdown()
		default:
			color.Red("[goTES3MP]: " + "Command" + ` "` + command + `" ` + "was not recognised.")
		}
	}
}

func initLogger() {
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
