package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

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
