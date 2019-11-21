package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// LoadConfig loads json config file
func LoadConfig() (ConfigLoaded bool) {
	var configPath = "./goTes3mp_config.json"
	viper.SetConfigName("goTes3mp_config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("tes3mpPath", ".")
	viper.SetDefault("debug", false)
	viper.SetDefault("enableCommands", true)
	viper.SetDefault("discordToken", "")
	viper.SetDefault("serverName", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs(configPath)
			fmt.Println(tes3mpLogMesage + "Created default config")
		} else {
			panic(fmt.Errorf(tes3mpLogMesage+"Fatal error reading config file: %s \n", err))
		}
	}
	fmt.Println(tes3mpLogMesage + "Successfully loaded config")

	return ConfigLoaded
}
