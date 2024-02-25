local cjson = require("cjson")
local goTES3MPUtils = require("custom.goTES3MP.utils")

local serverID = goTES3MP.GetServerID()
local discordServer = goTES3MP.GetDefaultDiscordServer()

local getJournal = {}

--- Get the journal entries for a player and send them as a message to a Discord channel
---@param playerName string The name of the player
---@param questID string The ID of the quest
---@param discordReplyChannel string The Discord channel to send the journal entries message to
---@return string The response message
getJournal.GetJournalEntries = function(playerName, questID)
    local questIndexs = {}

    -- Get the player by name
    local player = logicHandler.GetPlayerByName(playerName)
    if player then
        -- Iterate over each quest in the player's journal
        for _, quest in pairs(player.data.journal) do
            -- Check if the quest ID matches the provided quest ID (case-insensitive)
            if string.lower(quest.quest) == string.lower(questID) then
                table.insert(questIndexs, quest.index)
            end
        end

        -- Check if no matching quest entries were found
        if #questIndexs == 0 then
            return "Quest ID is invalid or player does not have this Quest."
        end

        -- Sort the quest indices in alphanumeric order
        questIndexs = goTES3MPUtils.alphanumsort(questIndexs)

        local questList = {}
        questList[#questList + 1] = "**" .. playerName .. "'s Journal entries for " .. '"' .. string.lower(questID) .. '"' .. "**\n"
        questList[#questList + 1] = "```"

        -- Concatenate the quest indices into a string
        for i, index in pairs(questIndexs) do
            questList[#questList + 1] = index
            if i < #questIndexs then
                questList[#questList + 1] = ","
            end
        end

        questList[#questList + 1] = "```"

        -- Send the quest list as a message to the Discord channel
        return table.concat(questList)
    else
        -- Send a message to the Discord channel indicating that the player does not exist
        return "Player does not exist."
    end
end

return getJournal