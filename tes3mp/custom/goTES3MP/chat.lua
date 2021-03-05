local cjson = require("cjson")
local goTES3MPChat = {}
local DiscordChannel = ""
local DiscordServer = ""

goTES3MPChat.ConfigureDiscord = function(discordserver, discordchannel)
    DiscordServer = DiscordServer
    DiscordChannel = DiscordChannel
end

customEventHooks.registerHandler("OnPlayerAuthentified", function(eventStatus, pid)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = GOTES3MPServerID,
        syncid = GoTES3MPSyncID,
        data = {
            channel = GoTES3MP_DiscordChannel,
			server = GoTES3MP_DiscordServer,
			message = "**"..tes3mp.GetName(pid) .. " has connected".."**"
        }
    }

    IrcBridge.SendSystemMessage(cjson.encode(messageJson))
end)

customEventHooks.registerValidator("OnPlayerDisconnect", function(eventStatus, pid)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = GOTES3MPServerID,
        syncid = GoTES3MPSyncID,
        data = {
            channel = GoTES3MP_DiscordChannel,
			server = GoTES3MP_DiscordServer,
			message = "**"..tes3mp.GetName(pid) .. " has disconnected".."**"
        }
    }
    IrcBridge.SendSystemMessage(cjson.encode(messageJson))
end)

customEventHooks.registerValidator(
    "OnPlayerSendMessage",
    function(eventStatus, pid, message)
        local messageJson = {
            method = "rawDiscord",
            source = "TES3MP",
            serverid = GOTES3MPServerID,
            syncid = GoTES3MPSyncID,
            data = {
                channel = GoTES3MP_DiscordChannel,
                server = GoTES3MP_DiscordServer,
                message = tes3mp.GetName(pid) ..": ".. message
            }
        }

        if message:sub(1, 1) == "/" then
		    return
        else
            responce = cjson.encode(messageJson)
            IrcBridge.SendSystemMessage(responce)
        end
    end
)
return goTES3MPChat