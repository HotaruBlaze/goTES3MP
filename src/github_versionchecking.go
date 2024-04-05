package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	color "github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

// GitHub file content response structure
type GitHubContentResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// Fetch file content from GitHub repository
func fetchIrcBridgeVersionFromGithub() (string, bool) {
	var err error
	// Construct GitHub API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", "hotarublaze", "gotes3mp", "tes3mp/scripts/custom/IrcBridge/IrcBridge.lua")

	// Send HTTP GET request to fetch file content
	resp, err := http.Get(url)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}

	// Parse GitHub response
	var contentResponse GitHubContentResponse
	if err := json.Unmarshal(body, &contentResponse); err != nil {
		return "", false
	}

	// Decode content from base64
	content, err := decodeBase64String(contentResponse.Content)
	if err != nil {
		return "", false
	}

	// Parse version number
	githubIrcBridgeVersion, err := parseVersionNumber(content)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	message, safeUpdate := compareVersions(githubIrcBridgeVersion)

	return message, safeUpdate
}

// Decode base64 encoded string
func decodeBase64String(encodedString string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

// Parse version number from file content
func parseVersionNumber(content string) (string, error) {
	// Regular expression to match version number
	re := regexp.MustCompile(`IrcBridge\.version\s*=\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("version number not found")
	}
	return matches[1], nil
}

// This compares the current version to the version in the github repo
// this is designed to alert the user if they can safely update the bot or not.
func compareVersions(githubVersion string) (string, bool) {
	versionCheck := strings.Compare(githubVersion, ircBridgeVersion)

	switch {
	// Github version is older
	case versionCheck < 0:
		// This shouldnt happen unless it's a custom build.
		return "Github version is newer than this code was built for, this usually only happens for custom builds", false
	case versionCheck > 0:
		return "Update is available however requires updating tes3mp's lua files", false
	default:
		return "No Lua Updates required", true
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
		// Check if the Lua updates are safe
		_, isSafeLuaUpdate := fetchIrcBridgeVersionFromGithub()
		if isSafeLuaUpdate {
			// If Lua updates are safe
			color.HiGreen("A new version of goTES3MP is available:  " + Build + " -> " + UpdateVersion + ".\nYour IRC bridge version should be compatible.")
		} else {
			// If Lua updates are not safe
			color.HiYellow("A new version of goTES3MP is available: " + Build + " -> " + UpdateVersion + ".\nHowever, the IRC bridge version may not be compatible.")
		}
	}

	if firstLaunch {
		color.HiBlack(strings.Repeat("=", 32))
	}
}

func checkforUpdates() {
	isUpdate, UpdateVersion := getLatestGithubRelease()

	if isUpdate {
		// Check if the Lua updates are safe
		_, isSafeLuaUpdate := fetchIrcBridgeVersionFromGithub()
		if isSafeLuaUpdate {
			// If Lua updates are safe
			color.HiGreen("A new version of goTES3MP is available:  " + Build + " -> " + UpdateVersion + ".\nYour IRC bridge version should be compatible.")
		} else {
			// If Lua updates are not safe
			color.HiYellow("A new version of goTES3MP is available: " + Build + " -> " + UpdateVersion + ".\nHowever, the IRC bridge version may not be compatible.")
		}
	}
}
