package main

import (
	"fmt"
	"math/big"
	"regexp"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// AppendIfMissing : Appends string if missing from array.
func AppendIfMissing(slice []string, i string) []string {
	currentSlice := slice
	for _, ele := range currentSlice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// // RemoveEntryFromArray : Remove Entry from Array.
// func RemoveEntryFromArray(array []string, remove string) []string {
// 	workArr := array
// 	for i := 0; i < len(workArr); i++ {
// 		if workArr[i] == remove {
// 			workArr = append(workArr[:i], workArr[i+1:]...)
// 			i--
// 		}
// 	}
// 	return workArr
// }

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%06X", n) // or %X or upper case
}

// FindinArray : Search String array for a value
func FindinArray(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// stringVerifier: Verify the string is "Clean"
// func stringVerifier(removeRGB bool, str string) string {
// 	message := str

// 	if removeRGB {
// 		message = removeRGBHex(message)
// 	}
// 	return message
// }

// removeRGBHex: Remove all RGB Hex's from string
func removeRGBHex(s string) string {
	message := s
	regex := "(?i)#[0-9A-F]{6}|#[0-9A-F]{3}"
	re := regexp.MustCompile(regex)

	message = re.ReplaceAllString(message, "")
	return message
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Infof("Alloc = %v MiB, TotalAlloc = %v MiB, Sys = %v MiB, NumGC = %v",
		bToMb(m.Alloc),
		bToMb(m.TotalAlloc),
		bToMb(m.Sys),
		m.NumGC,
	)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// convertDiscordFormattedItems : Formats Discord Formatted Items to Text Only
func convertDiscordFormattedItems(str string, gid string) string {
	// match all discord formatted items
	reChannel := regexp.MustCompile(`<#(\d+)>`)
	reMentionUser := regexp.MustCompile(`<@(?:\!)?(\d+)>`)
	reMentionRole := regexp.MustCompile(`<@&(\d+)>`)

	str = reChannel.ReplaceAllStringFunc(str, func(s string) string {
		// get the channels for this guild
		id := reChannel.FindStringSubmatch(s)[1]
		channels, err := DiscordSession.GuildChannels(gid)
		if err != nil {
			log.Errorln("[utils:convertDiscordFormattedItems]", "Error getting guild channels: ", err)
			return s
		}
		for _, c := range channels {
			if c.ID == id {
				return "%" + c.Name
			}
		}

		return s
	})
	reMentionUser.FindAllStringSubmatchIndex(str, -1)
	str = reMentionUser.ReplaceAllStringFunc(str, func(s string) string {
		// get the user for this guild
		id := reMentionUser.FindStringSubmatch(s)[1]
		member, err := DiscordSession.GuildMember(gid, id)
		if err != nil {
			log.Errorln("[utils:convertDiscordFormattedItems]", "Error getting guild member: ", err)
			return s
		}
		if member.Nick != "" {
			return member.Nick
		} else {
			return member.User.Username
		}
	})
	str = reMentionRole.ReplaceAllStringFunc(str, func(s string) string {
		// get the role for this guild
		id := reMentionRole.FindStringSubmatch(s)[1]
		roles, err := DiscordSession.GuildRoles(gid)
		if err != nil {
			log.Errorln("[utils:convertDiscordFormattedItems]", "Error getting guild roles: ", err)
			return s
		}
		for _, r := range roles {
			if r.ID == id {
				return "@" + r.Name
			}
		}
		return s
	})
	return str
}

// filterDiscordEmotes : Formats Discord Emotes Correctly
func filterDiscordEmotes(str string) string {
	re := regexp.MustCompile(`<:(\S+):\d+>`)
	filteredString := re.ReplaceAllString(str, `:$1:`)

	return filteredString
}

// MemoryDebugInfo : Print current memory and GC cycles, Used for monitoring for memory leaks
func MemoryDebugInfo() {
	printMemUsage()
	for {
		time.Sleep(30 * time.Minute)
		printMemUsage()
	}

}

// func trimLeftChar(s string) string {
// 	for i := range s {
// 		if i > 0 {
// 			return s[i:]
// 		}
// 	}
// 	return s[:0]
// }

func checkError(str string, err error) bool {
	if err != nil {
		log.Errorf(str, err)
		return true
	}
	return false
}
