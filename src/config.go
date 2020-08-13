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
	viper.SetDefault("webserver.port", ":8080")

	viper.SetDefault("tes3mp.baseDir", ".")
	viper.SetDefault("debug", false)
	viper.SetDefault("serveroutput", true)
	viper.SetDefault("commandPrefix", "!")
	viper.SetDefault("printMemoryInfo", false)

	viper.SetDefault("irc.enableChatChannel", false)
	viper.SetDefault("irc.server", "127.0.0.1")
	viper.SetDefault("irc.port", "6667")
	viper.SetDefault("irc.nick", "goTES3MP")
	viper.SetDefault("irc.systemchannel", "#goTES3MP-System")
	viper.SetDefault("irc.chatchannel", "#goTES3MP-Chat")
	viper.SetDefault("irc.pass", "")

	viper.SetDefault("discord.enableCommands", true)
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
