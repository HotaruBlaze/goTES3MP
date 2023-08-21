local getJournal = {}
local cjson = require("cjson")
local goTES3MPUtils = require("custom.goTES3MP.utils")

getJournal.GetJournalEntrys = function(playerName, questID, discordReplyChannel)
    local questIndexs = {}
    local player = logicHandler.GetPlayerByName(playerName)
    if player ~= nil then
        for _, quest in pairs(player.data.journal) do
            if string.lower(quest.quest) == string.lower(questID) then
                table.insert(questIndexs, quest.index)
            end
        end
        if #questIndexs == 0 then
            goTES3MPUtils.sendDiscordMessage(
                goTES3MP.GetServerID(),
                discordReplyChannel,
                goTES3MP.GetDefaultDiscordServer(),
                "**Quest ID is invalid or player does not have this Quest.**"
            )
            return
        end

        questIndexs = goTES3MPUtils.alphanumsort(questIndexs)
        questList =
            "**" .. playerName .. "'s Journal entry's for " .. '"' .. string.lower(questID) .. '"' .. "**" .. "\n"
        questList = questList .. "```" .. "\n"

        for i, index in pairs(questIndexs) do
            if i == #questIndexs then
                questList = questList .. index .. "\n"
            else
                questList = questList .. index .. ","
            end
        end
        questList = questList .. "```"

        goTES3MPUtils.sendDiscordMessage(
            goTES3MP.GetServerID(),
            discordReplyChannel,
            goTES3MP.GetDefaultDiscordServer(),
            questList
        )
    else
        goTES3MPUtils.sendDiscordMessage(
            goTES3MP.GetServerID(),
            discordReplyChannel,
            goTES3MP.GetDefaultDiscordServer(),
            "**"  .. "Player does not Exist." .. "**"
        )
    end
end

return getJournal