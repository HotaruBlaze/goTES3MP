local getPlayers = {}
local cjson = require("cjson")
local goTES3MPModules = goTES3MP.GetModules()


getPlayers.getPlayers = function(discordReplyChannel)
    local playerList = ""

    if tableHelper.getCount(Players) > 0 then
        for pid, player in pairs(Players) do
            if player ~= nil and player:IsLoggedIn()then
                playerList = playerList .. player.name .. "\n"
            end
        end
    end

    -- Check if playerList is empty or has no players
    if playerList == "" then
        local noPlayersMessage = "**No players are currently online.**"
        goTES3MPModules["utils"].sendDiscordMessage(
            goTES3MP.GetServerID(),
            discordReplyChannel,
            goTES3MP.GetDefaultDiscordServer(),
            noPlayersMessage
        )
    else
        playerList = "```" .."\n".. playerList .."\n".. "```"

        goTES3MPModules["utils"].sendDiscordMessage(
            goTES3MP.GetServerID(),
            discordReplyChannel,
            goTES3MP.GetDefaultDiscordServer(),
            playerList
        )
    end
end

return getPlayers