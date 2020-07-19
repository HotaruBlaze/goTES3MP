package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GitCommit used to build
var GitCommit string

var Stdin io.WriteCloser

// Build version
var Build = "v0.0.0-Debug"

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

// TES3MPBuild : Linux/Windows (32-bit/64-Bit)
var TES3MPBuild = ""

var tes3mpLogMessage = "[goTES3MP]"

// MultiWrite : Prints to logfile and os.Stdout
var MultiWrite io.Writer

func main() {
	queue := goconcurrentqueue.NewFIFO()
	printBuildInfo()
	initLogger()
	LoadConfig()
	enableDebug := viper.GetBool("debug")
	if enableDebug {
		log.Warnln("Debug mode is enabled")
		log.SetLevel(log.DebugLevel)
	}
	if viper.GetBool("webserver.enable") {
		go InitWebserver()
	}
	go UpdateStatusTimer()
	if viper.GetBool("discord.enable") {
		go InitDiscord()
	}
	if viper.GetBool("irc.enable") {
		go InitIRC()
	}
	if viper.GetBool("printMemoryInfo") {
		go MemoryDebugInfo()
	}

	go queueProcessor(queue)
	LaunchTes3mp(queue)
}

func initLogger() {
	dt := time.Now()
	ProgramDirectory := "./goTES3MP/logs/"
	logfileName := ProgramDirectory + "goTES3MP-" + dt.Format("02-01-2006-15_04_05") + ".log"

	if _, err := os.Stat(ProgramDirectory); os.IsNotExist(err) {
		os.MkdirAll(ProgramDirectory, 0700)
	}

	logfile, err := os.OpenFile(logfileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	MultiWrite = io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(MultiWrite)
	if viper.GetBool("debug") {
		println("DEBUG IS ON")
	}
	log.SetLevel(log.InfoLevel)
	if Build != "" && GitCommit != "" {
		log.Infoln("goTES3MP", "Build:", Build+", "+"Commit:", GitCommit)
	}
}

func printBuildInfo() {
	fmt.Println("================================")
	fmt.Println("goTES3MP")
	fmt.Println("Build:", Build)
	fmt.Println("Commit:", GitCommit)
	fmt.Println("Github:", "https://github.com/hotarublaze/goTES3MP")
	fmt.Println("================================")

}

// LaunchTes3mp : Start and initialize TES3MP
func LaunchTes3mp(queue *goconcurrentqueue.FIFO) {
	tes3mpPath := viper.GetString("tes3mp.basedir")

	tes3mpBinary := "/tes3mp-server"

	cmd := exec.Command(tes3mpPath + tes3mpBinary)

	Stdin, _ = cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	startErr := cmd.Start()
	if startErr != nil {
		log.Fatalf("cmd.Run() failed with '%s'\n", startErr)
	}

	//// Does not seem to shutdown tes3mp correctly
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Infoln("Recieved Signal to exit, Exiting and notifying discord")
		Stdin.Write([]byte("\n"))
	}()
	////

	outScanner := bufio.NewScanner(stdout)
	for outScanner.Scan() {
		queue.Enqueue(outScanner.Text())
	}
}
func queueProcessor(queue *goconcurrentqueue.FIFO) {
	for {
		str, err := queue.DequeueOrWaitForNextElement()
		if err != nil {
			fmt.Println(err)
			return
		}
		tes3mpOutputHandler(str.(string))
	}
}
