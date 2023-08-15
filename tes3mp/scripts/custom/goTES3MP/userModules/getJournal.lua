local cjson = require("cjson")
local goTES3MPModules = goTES3MP.GetModules()

local getJournal = {}

getJournal.GetJournalEntries = function(playerName, questID, discordReplyChannel)
    local player = logicHandler.GetPlayerByName(playerName)

    if player == nil then
        goTES3MPModules["utils"].sendDiscordMessage(
            goTES3MP.GetServerID(),
            discordReplyChannel,
            goTES3MP.GetDefaultDiscordServer(),
            "**Player does not exist.**"
        )
        return
    end

    local questIndexes = {}
    for _, quest in pairs(player.data.journal) do
        if string.lower(quest.quest) == string.lower(questID) then
            table.insert(questIndexes, quest.index)
        end
    end

    if #questIndexes == 0 then
        goTES3MPModules["utils"].sendDiscordMessage(
            goTES3MP.GetServerID(),
            discordReplyChannel,
            goTES3MP.GetDefaultDiscordServer(),
            "**Quest ID is invalid or player does not have this Quest.**"
        )
        return
    end

    table.sort(questIndexes, goTES3MPModules["utils"].alphanumsort)

    local questList = {
        "**" .. playerName .. "'s Journal entries for " .. '"' .. string.lower(questID) .. '"' .. "**\n",
        "```\n"
    }

    for i, index in ipairs(questIndexes) do
        if i == #questIndexes then
            table.insert(questList, index .. "\n")
        else
            table.insert(questList, index .. ",")
        end
    end

    table.insert(questList, "```")

    goTES3MPModules["utils"].sendDiscordMessage(
        goTES3MP.GetServerID(),
        discordReplyChannel,
        goTES3MP.GetDefaultDiscordServer(),
        table.concat(questList)
    )
end

return getJournal
