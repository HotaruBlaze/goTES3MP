# goTES3MP for TES3MP

An golang Client/Server appliaction for TES3MP utilizing golang and IRC for server communication and features such as a Discord-Chat bridge, anti-VPN protection and more. 

The previous depricated version of this is: [TES3MP_DiscordRelay](https://github.com/HotaruBlaze/TES3MP_DiscordRelay)

# Features
- Discord-Chatbridge to/from TES3MP and Discord
- Anti-VPN protection, using publicly available apis
- Ability to do remote TES3MP commands from Discord.
# Worthwhile Notes
* TES3MP does not shutdown correctly most of the time with SIGINT, or closing the application, It's recommended to use another script for this, such as [ShutdownServer](https://github.com/tes3mp-scripts/ShutdownServer).

* Sometimes goTES3MP fails to connect to Discord, this usually fixes itself after a couple of minutes, if not try running the "reloaddiscord" command on goTES3MP. 

# Requirements
- Golang version >= 1.20 
- An IRC Server, I recommend my personal fork of [oragono](https://github.com/oragono/oragono) found [here](https://github.com/HotaruBlaze/oragono-tes3mp)
- [Datamanager](https://github.com/tes3mp-scripts/DataManager) for TES3MP
- *[cjson](https://github.com/TES3MP/lua-cjson) (Included in tes3mp-scripts.zip)

# Install Instructions - Standalone
1. Download the latest build with accompanying tes3mp-scripts.zip 
2. Extract and copy the custom and lib folders to `server` folder.
3. Add the following to your server/customScripts.lua file, making sure DataManager is above the following
```
IrcBridge = require("custom/IrcBridge/IrcBridge")
goTES3MP = require("custom/goTES3MP/main")
```
4. Download and place the correct `goTES3MP` binary for your platform
5. Run the binary to generate the default configuration file(`config.yaml`)

# Install Instructions - Docker-Compose
```yml
version: "3"
services:
  irc-server:
    image: mrflutters/oragono:tes3mp-fork
    ports:
      - 172.17.0.1:6667:6667 #Plaintext
    restart: unless-stopped
    volumes:
        - irc_data:/ircd
        - ./oragono/ircd.yaml:/ircd/ircd.yaml
        - ./oragono/fullchain.pem:/ircd/fullchain.pem
        - ./oragono/privkey.pem:/ircd/privkey.pem
    networks:
      - relay-net
    container_name: irc-server

  gotes3mp:
    image: 'ghcr.io/hotarublaze/gotes3mp:v0.3.4'
    volumes:
      - './logs:/app/logs'
      - './config.yaml:/app/config.yaml'
    networks:
      - relay-net
      
networks:
  relay-net:
volumes:
  irc_data:
```
# Build Instructions - Linux (assumes Ubuntu/Debian)
```
sudo apt install golang-go git
git clone https://github.com/HotaruBlaze/goTES3MP
cd goTES3MP
./scripts/build.sh
```