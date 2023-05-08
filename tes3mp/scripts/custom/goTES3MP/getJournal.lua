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
            local messageJson = {
                method = "rawDiscord",
                source = "TES3MP",
                serverid = goTES3MP.GetServerID(),
                syncid = GoTES3MPSyncID,
                data = {
                    channel = discordReplyChannel,
                    server = goTES3MP.GetDefaultDiscordServer(),
                    message = "**Quest ID is invalid or player does not have this Quest.**"
                }
            }
            local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
            if responce ~= nil then
                IrcBridge.SendSystemMessage(responce)
            end
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

        local messageJson = {
            method = "rawDiscord",
            source = "TES3MP",
            serverid = goTES3MP.GetServerID(),
            syncid = GoTES3MPSyncID,
            data = {
                channel = discordReplyChannel,
                server = goTES3MP.GetDefaultDiscordServer(),
                message = questList
            }
        }

        local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
        if responce ~= nil then
            IrcBridge.SendSystemMessage(responce)
        end
    else
        local messageJson = {
            method = "rawDiscord",
            source = "TES3MP",
            serverid = goTES3MP.GetServerID(),
            syncid = GoTES3MPSyncID,
            data = {
                channel = discordReplyChannel,
                server = goTES3MP.GetDefaultDiscordServer(),
                message = "**"  .. "Player does not Exist." .. "**"
            }
        }

        local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
        if responce ~= nil then
            IrcBridge.SendSystemMessage(responce)
        end
    end
end

return getJournal
