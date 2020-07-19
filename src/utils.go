package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// AppendIfMissing : Appends string if missing from array.
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// RemoveEntryFromArray : Remove Entry from Array.
func RemoveEntryFromArray(array []string, remove string) []string {
	workArr := array
	for i := 0; i < len(workArr); i++ {
		if workArr[i] == remove {
			workArr = append(workArr[:i], workArr[i+1:]...)
			i-- // form the remove item index to start iterate next item
		}
	}
	return workArr
}
func removeString(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

// AppConfigProperties String map for reading tes3mp-server-default.cfg
type AppConfigProperties map[string]string

// ReadPropertiesFile for reading .cfg files to correctly read its values
func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	TES3MPServerConfig := AppConfigProperties{}

	if len(filename) == 0 {
		return TES3MPServerConfig, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				TES3MPServerConfig[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return TES3MPServerConfig, nil
}
func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
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
func stringVerifier(removeRGB bool, str string) string {
	message := str

	if removeRGB {
		message = removeRGBHex(message)
	}
	return message
}

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

// MemoryDebugInfo : Print current memory and GC cycles, Used for monitoring for memory leaks
func MemoryDebugInfo() {
	printMemUsage()
	for {
		time.Sleep(30 * time.Minute)
		printMemUsage()
	}

}
