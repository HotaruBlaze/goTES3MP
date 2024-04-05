//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	// Read Lua file content
	luaContent, err := os.ReadFile("../tes3mp/scripts/custom/IrcBridge/IrcBridge.lua")
	if err != nil {
		fmt.Println("Error reading Lua file:", err)
		return
	}

	// Extract version using regular expression
	re := regexp.MustCompile(`IrcBridge\.version\s*=\s*"([^"]+)"`)
	matches := re.FindSubmatch(luaContent)
	if len(matches) < 2 {
		fmt.Println("Version not found in Lua file")
		return
	}

	version := string(matches[1])

	// Write the version to version.go
	fileContent := fmt.Sprintf(`package main

var ircBridgeVersion = "%s"
`, version)

	// Write the generated Go code to version.go
	if err := os.WriteFile("version.go", []byte(fileContent), 0644); err != nil {
		fmt.Println("Error writing version:", err)
		return
	}

	fmt.Println("Version set to", version)
}
