# goTES3MP for TES3MP

This is an attempt at creating a bot for tes3mp-server designed for running more intensive operations with minimal lua.

# Known Issues
* TES3MP does not shutdown correctly most of the time with SIGINT, or closing the application, It's recommended to use another script for this, such as [ShutdownServer](https://github.com/tes3mp-scripts/ShutdownServer).
* Sometimes server output from TES3MP will no longer be printed into the terminal, The reason is currently unclear.
* Sometimes goTES3MP loses connection to Discord, this usually fixes itself after a couple of minutes, if not try restarting goTES3MP. 

# Goal
- [X] Create a extendable application for manipulating and handling TES3MP output
- [x] Recreate [TES3MP_DiscordRelay](https://github.com/HotaruBlaze/TES3MP_DiscordRelay) with bug fixes and Discord role support.
- [X] Added a web endpoint for accessing server status, such as current player count and players.
- [X] Show CurrentPlayers/MaxPlayers as Discord bot status.

# Requirements
- Linux 
- An IRC Server, I recommend [oragono](https://github.com/oragono/oragono)
- [Datamanager](https://github.com/tes3mp-scripts/DataManager) for TES3MP
- [cjson](https://github.com/TES3MP/lua-cjson) (**Not Required but highly recommended and included in tes3mp-scripts.zip.**)
# Install Instructions
1. Download the latest build with accompanying tes3mp-scripts.zip 
2. Extract and copy the custom and lib folders to `server` folder.
3. Add the following to your server/customScripts.lua file, making sure DataManager is above the following
```
IrcBridge = require("custom/IrcBridge/IrcBridge")
goTES3MP = require("custom/goTES3MP/main")
```
4. Place `goTES3MP-Linux` in the same directory as TES3MP
5. Run `goTES3MP-Linux` to generate the default configuration file(`goTes3mp_config.json`)
