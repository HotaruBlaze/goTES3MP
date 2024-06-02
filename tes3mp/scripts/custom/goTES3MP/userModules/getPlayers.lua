local getPlayers = {}
local goTES3MPModules = goTES3MP.GetModules()


--- Retrieve the list of players and send it as a message to a Discord channel
---@param discordReplyChannel string The Discord channel to send the player list message to
---@return string The player list message
getPlayers.getPlayers = function()
    local playerList = ""

    -- Check if there are any players online
    if tableHelper.getCount(Players) > 0 then
        -- Iterate over each player
        for pid, player in pairs(Players) do
            -- Check if the player is logged in
            if player ~= nil and player:IsLoggedIn() then
                -- Add the player's name to the player list
                playerList = playerList .. player.name .. "\n"
            end
        end
    end

    -- Check if playerList is empty or has no players
    if playerList == "" then
        local noPlayersMessage = "**No players are currently online.**"
        return noPlayersMessage

    else
        -- Format the playerList with triple backticks and send it as a message to Discord
        playerList = "```" .."\n".. playerList .."\n".. "```"
        return playerList
    end
end

return getPlayers