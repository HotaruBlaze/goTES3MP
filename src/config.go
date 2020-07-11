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
	viper.AddConfigPath("./goTES3MP")
	viper.AddConfigPath(".")

	viper.SetDefault("webserver.enable", false)

	viper.SetDefault("tes3mp.baseDir", ".")
	viper.SetDefault("debug", false)
	viper.SetDefault("enable_ServerOutput", true)
	viper.SetDefault("commandPrefix", "!")

	viper.SetDefault("irc.enable", false)
	viper.SetDefault("irc.server", "")
	viper.SetDefault("irc.port", "")
	viper.SetDefault("irc.nick", "goTES3MP")
	viper.SetDefault("irc.systemchannel", "#")
	viper.SetDefault("irc.chatchannel", "#")
	viper.SetDefault("irc.pass", "6667")

	// viper.SetDefault("enableCommands", true)
	viper.SetDefault("discord.enable", false)
	viper.SetDefault("discord.allowColorHexUsage", false)
	viper.SetDefault("discord.token", "")
	viper.SetDefault("discord.alertsChannel", "")
	viper.SetDefault("discord.serverChat", "")
	viper.SetDefault("discord.staffRoles", []string{})
	viper.SetDefault("discord.userRoles", []string{})

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
