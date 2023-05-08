local getPlayers = {}
local cjson = require("cjson")
local goTES3MPUtils = require("custom.goTES3MP.utils")

getPlayers.getPlayers = function(discordReplyChannel)
    local playerList = ""

    if tableHelper.getCount(Players) > 0 then
        for pid, player in pairs(Players) do
            if player ~= nil and player:IsLoggedIn()then
                playerList = playerList .. player.name .. "\n"
            end
        end
    end

    playerList = "```" .."\n".. playerList .."\n".. "```"

    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = goTES3MP.GetServerID(),
        syncid = GoTES3MPSyncID,
        data = {
            channel = discordReplyChannel,
            server = goTES3MP.GetDefaultDiscordServer(),
            message = playerList
        }
    }

    local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
    if responce ~= nil then
        IrcBridge.SendSystemMessage(responce)
    end
end

return getPlayers