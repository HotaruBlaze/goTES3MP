# goTES3MP for TES3MP
### Note:
This was an attempt at having an aditional layer ontop of tes3mp, however this proved to be unreliable. So this was replaced with a client-server model, that allows you to move more advanced logic to golang, And by default serves as a replacement for [TES3MP_DiscordRelay](https://github.com/HotaruBlaze/TES3MP_DiscordRelay)

# Known Issues
<!-- * TES3MP does not shutdown correctly most of the time with SIGINT, or closing the application, It's recommended to use another script for this, such as [ShutdownServer](https://github.com/tes3mp-scripts/ShutdownServer). -->
<!-- * Sometimes goTES3MP loses connection to Discord, this usually fixes itself after a couple of minutes, if not try running the "reloaddiscord" command.  -->
* Mentioning Channel names and emotes will be glitched/formatted incorrectly, This is known and a fix is being looked for.

# Goal
- [x] Recreate [TES3MP_DiscordRelay](https://github.com/HotaruBlaze/TES3MP_DiscordRelay) with bug fixes and Discord role support.
- [X] Added a web endpoint for accessing server status, such as current player count and players.
- [X] Show CurrentPlayers/MaxPlayers as Discord bot status.

# Requirements
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
    image: 'ghcr.io/hotarublaze/gotes3mp:v0.3'
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