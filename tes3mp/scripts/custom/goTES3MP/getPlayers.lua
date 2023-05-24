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

    goTES3MPUtils.sendDiscordMessage(
        goTES3MP.GetServerID(),
        discordReplyChannel,
        goTES3MP.GetDefaultDiscordServer(),
        playerList
    )
end

return getPlayers