package main

import (
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfig loads json config file
func LoadConfig() (ConfigLoaded bool) {
	var configPath = "./config.yaml"
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("webserver.enable", false)
	viper.SetDefault("webserver.port", ":8080")

	viper.SetDefault("tes3mp.serverid", "")
	viper.SetDefault("debug", false)
	viper.SetDefault("discord.commandPrefix", "!")
	viper.SetDefault("printMemoryInfo", false)

	viper.SetDefault("irc.enableChatChannel", false)
	viper.SetDefault("irc.server", "127.0.0.1")
	viper.SetDefault("irc.port", "6667")
	viper.SetDefault("irc.nick", "goTES3MP")
	viper.SetDefault("irc.systemchannel", "#goTES3MP-System")
	viper.SetDefault("irc.chatchannel", "#goTES3MP-Chat")
	viper.SetDefault("irc.pass", "")

	viper.SetDefault("vpn.iphub_apikey", "")
	viper.SetDefault("vpn.ipqualityscore_apikey", "")

	viper.SetDefault("discord.enableCommands", true)
	viper.SetDefault("discord.boldPlayerNames", false)
	viper.SetDefault("discord.enable", false)
	viper.SetDefault("discord.allowColorHexUsage", false)
	viper.SetDefault("discord.token", "")
	viper.SetDefault("discord.alertsChannel", "")
	viper.SetDefault("discord.serverChat", "")
	viper.SetDefault("discord.staffRoles", []string{})
	viper.SetDefault("discord.userroles", []string{})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err := viper.WriteConfigAs(configPath)
			if err != nil {
				log.Errorln("[Config]", "Failed to write Config: ", err)
			}
			log.Infoln("[Viper]", "Created default config")
			os.Exit(1)
		} else {
			log.Errorf("[Viper]", "Fatal error reading config file: %v", err)
			panic(1)
		}
	}
	log.Info("[Viper]", "Successfully loaded config")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("[Viper] Reloaded Configuration file")
	})
	return ConfigLoaded
}
