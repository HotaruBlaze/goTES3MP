# GoTES3MP Methods

### Note: syncid is not used however it hasnt been fully removed from the source code and still has to be included in some places, however can be a blank string.

# "Sync" Method: [Link](../tes3mp/scripts/custom/goTES3MP\sync.lua)
```lua
local messageJson = {
    ServerID = ServerID,
    method = "Sync",
    source = "TES3MP",
    data = {
        MaxPlayers = tostring(tes3mp.GetMaxPlayers()),
        CurrentPlayerCount = tostring(logicHandler.GetConnectedPlayerCount()),
        Forced = tostring(forceResync),
        Status = "Ping",
    }
}
```
The Sync method is very importent to goTES3MP, this is designed to make sure the communications are working correctly, aswell as keeping the player count updated.

You may also threat this as a Ping/Pong system, as goTES3MP will respond with a pong if it reads and is able to process what it received. 

The **Forced** variable is used with the command `/forceSync` and will force goTES3MP to accept the data thats being sent with it, this is useful if your ServerID has changed or you restarted/updated goTES3MP, however this is usually not needed and will be handled automatically.

# "IRC" Method:
Thie script used to support using IRC to **<u>Chat with Discord and tes3mp</u>**. However when this was rewritten with a more refined system, this method was never reimplemented, It's been left in as it may be reimplemented in the future.

# "rawDiscord" Method:
```lua
local messageJson = 
    method = "rawDiscord",
    source = "TES3MP",
    serverid = goTES3MP.GetServerID(),
    syncid = GoTES3MPSyncID,
    data = {
        channel = discordChannel,
        server = discordServer,
        message = "**" .. "[TES3MP] " .. tes3mp.GetName(pid) .. " has connected" .. "**"
    }
```
This method is the main one you will be using for tes3mp->discord, as this method is very flexible with its use. 

channel: The desired discord channel you want the bot to send the message to.
server: The desired discord server you want the bot to use.

message: This is what you want the bot to send to discord as a message. Note that this also supports multilined inputs and any discord formatting tips. This should also work with emotes however they can be a little buggy, so it's **<u>recommended that you use Discord emoticons and not default unicode, such as smiley faces</u>**


# "VPNCheck" Method: [Link](../tes3mp/scripts/custom/goTES3MP/VPNChecker.lua)
```lua
local messageJson = {
    method = "VPNCheck",
    source = "TES3MP",
    serverid = goTES3MP.GetServerID(),
    syncid = GoTES3MPSyncID,
    data = {
        channel = discordChannel,
        server = discordServer,
        message = IP,
        playerpid = tostring(pid)
    }
}
```
Now, this method is quite simple, all this method does is send goTES3MP the players PID and IP to goTES3MP with a VPNCheck method.

goTES3MP will then check that against multiple sources to see how trustworthy an ip is. Using the following websites:<br>
https://iphub.info<br>
https://ipqualityscore.com

You may find the related goTES3MP code [here](../src/vpnChecker.go), However theirs not much to read, it will just show you the api responces and how the bot is building/handling it.

Note that currently, if you wish to use Anti-VPN you must use both services, as iphub can miss some ip's that ipqualityscore will catch. <br><u>**Warning:** Their is currently no check in place if one of the api keys are missing and the application may misbehave.</u>

If an IP is deamed to be an VPN or seemingly bad, such as a proxy, it will send a VPNCheck method back to tes3mp, with a modified data packet, telling it to kick that pid. You can find this [Here](../tes3mp/scripts/custom/IrcBridge/IrcBridge.lua#L107-L130)
