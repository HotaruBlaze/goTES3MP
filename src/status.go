package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	color "github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
)

// ServerStatus struct
type ServerStatus struct {
	ServerOnline   bool
	TES3MPVersion  string
	Build          string
	Players        []string
	CurrentPlayers int
	MaxPlayers     int
}

// UpdateStatusTimer run UpdateStatus loop
func UpdateStatusTimer() {
	getServerMaxPlayers()
	for {
		time.Sleep(20 * time.Second)
		_ = UpdateStatus()
	}

}

func getServerMaxPlayers() {
	MaxPlayers = 0
}

// UpdateStatus for keeping server stats synced
func UpdateStatus() (s []byte) {
	CurrentStatus := ServerStatus{
		ServerOnline:   ServerRunning,
		TES3MPVersion:  TES3MPVersion,
		Build:          Build,
		Players:        Players,
		CurrentPlayers: CurrentPlayers,
		MaxPlayers:     MaxPlayers,
	}
	var jsonData []byte
	jsonData, err := json.Marshal(CurrentStatus)
	if err != nil {
		log.Println(err)
	}
	return jsonData
}

func getLatestGithubRelease() (isUpdate bool, latestVersion string) {
	client := github.NewClient(nil)
	releases, _, _ := client.Repositories.GetLatestRelease(context.Background(), "HotaruBlaze", "goTES3MP")
	latestRelease := releases.GetTagName()

	// Get the build number thats set on build
	currentBuild, err := version.NewVersion(Build)
	if err != nil {
		log.Println(err)
	}

	// Get latest github release
	latestBuild, err := version.NewVersion(latestRelease)
	if err != nil {
		log.Println(err)
	}

	if currentBuild.LessThan(latestBuild) {
		return true, string("v" + latestBuild.String())
	} else {
		return false, "nil"
	}

}

func getStatus(firstLaunch bool, showModules bool) {
	if firstLaunch {
		color.HiBlack(strings.Repeat("=", 32))
	}
	color.HiBlack("goTES3MP: " + Build)
	color.HiBlack("Commit: " + GitCommit)
	color.HiBlack("Github: " + "https://github.com/hotarublaze/goTES3MP" + "\n")
	color.HiBlack("Interactive Console: " + strconv.FormatBool(viper.GetBool("enableInteractiveConsole")))
	isUpdate, UpdateVersion := getLatestGithubRelease()
	if isUpdate {
		color.HiGreen("New build of goTES3MP is available: " + UpdateVersion + ", Current: " + Build)
	}

	if firstLaunch {
		color.HiBlack(strings.Repeat("=", 32))
	}
}
