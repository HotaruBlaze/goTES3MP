package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfig loads json config file
func LoadConfig() (ConfigLoaded bool) {
	var configPath = "./goTes3mp_config.json"
	viper.SetConfigName("goTes3mp_config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	viper.SetDefault("tes3mp.baseDir", ".")
	viper.SetDefault("discord.allowColorHexUsage", false)
	viper.SetDefault("debug", false)
	viper.SetDefault("enable_ServerOutput", true)
	viper.SetDefault("commandPrefix", "!")

	viper.SetDefault("irc.server", "")
	viper.SetDefault("irc.port", "")
	viper.SetDefault("irc.nick", "")
	viper.SetDefault("irc.channel", "")
	viper.SetDefault("irc.pass", "")

	// viper.SetDefault("enableCommands", true)
	viper.SetDefault("discord.token", "")
	viper.SetDefault("discord.serverChat", "")
	viper.SetDefault("discord.staffRoles", []string{})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfigAs(configPath)
			log.Infoln("[Viper]", "Created default config")
			os.Exit(1)
		} else {
			log.Errorf("[Viper]", "Fatal error reading config file: %s \n", err)
			panic(1)
		}
	}
	log.Println("[Viper]", "Successfully loaded config")

	return ConfigLoaded
}
