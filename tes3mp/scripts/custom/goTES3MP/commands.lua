local commands = {}
local goTES3MPModules = goTES3MP.GetModules()

---@class CommandHandler
---@field description string
---@field handler fun(player: string, commandArgs: string[], discordReplyChannel: string)
---@type table<string, CommandHandler>
local commandHandlers = {
    ["kickplayer"] = {
        description = "Kicks the specified player from the tes3mp server.",
        handler = function(playerPID, playerName, commandArgs, discordReplyChannel)
            targetPid = commands.getPlayerPID(table.concat(commandArgs, " "))
            if targetPid ~= nil then
                tes3mp.SendMessage(targetPid, color.Red .. "[SYSTEM] " .. "You have been kicked by an Administrator.", false)
                tes3mp.Kick(targetPid)
                commands.SendResponse(discordReplyChannel)
            end
        end
    },
    ["runconsole"] = {
        description = "Run a console command on a specific Player.",
        handler = function(playerPID, playerName, consoleCommand, discordReplyChannel)
            targetPid = commands.getPlayerPID(playerName)
            if targetPid ~= nil then
                logicHandler.RunConsoleCommandOnPlayer(targetPid, table.concat(consoleCommand, " "))
                commands.SendResponse(discordReplyChannel)
            end
        end
    },
    ["resetkills"] = {
        description = "Reset player kills.",
        handler = function(playerPID, playerName, commandArgs, discordReplyChannel)
            for refId, killCount in pairs(WorldInstance.data.kills) do
                WorldInstance.data.kills[refId] = 0
            end
            WorldInstance:QuicksaveToDrive()
            if tableHelper.getCount(Players) > 0 then
                for pid, player in pairs(Players) do
                    WorldInstance:LoadKills(pid, true)
                    tes3mp.SendMessage(pid, "All the kill counts for creatures and NPCs have been reset.\n", false)
                end
            end
            tes3mp.LogMessage(enumerations.log.INFO, "All the kill counts for creatures and NPCs have been reset.")
            commands.SendResponse(discordReplyChannel)
        end
    },
    ["players"] = {
        description = "List Players",
        handler = function(playerPID, playerName, commandArgs, discordReplyChannel)
            goTES3MPModules.getPlayers.getPlayers(discordReplyChannel)
        end
    },
    ["getjournal"] = {
        description = "get a player's Journal Entry",
        handler = function(playerPID, playerName, commandArgs, discordReplyChannel)
            local reconstructedString = table.concat(commandArgs, " ")
            local playerName, questID = reconstructedString:match('^"([^"]-)"%s*(.*)$')
            local playerPID = commands.getPlayerPID(playerName)
            goTES3MPModules.getJournal.GetJournalEntries(playerName, questID, discordReplyChannel)
        end
    },
}
-- Define the !help command handler separately
---@param player string
---@param commandArgs table
---@param discordReplyChannel string
local helpHandler = function(playerPID, playerName, commandArgs, discordReplyChannel)
    local commandList = "Available commands:\n"
    for cmd, data in pairs(commandHandlers) do
        commandList = commandList .. "!" .. cmd .. ": " .. data.description .. "\n"
    end
    commandList = "```\n" .. commandList .. "```"
    
    goTES3MPModules.utils.sendDiscordMessage(goTES3MP.GetServerID(), discordReplyChannel, goTES3MP.GetDefaultDiscordServer(), commandList)
end

-- Add the help handler to the commandHandlers table
commandHandlers["help"] = {
    description = "Display a list of available commands.",
    handler = helpHandler
}

--- Process the command and call the appropriate handler.
---@param command string
---@param commandArgs string[]
---@param discordReplyChannel string
commands.processCommand = function(command, commandArgs, discordReplyChannel)
    local command = string.lower(command)
    local commandHandlerData = commandHandlers[command]

    -- I hate this and it needs to be redesigned but i keep re-implementing this over and over with the same flaws.
    -- TODO: Refactor this, we can probably only do this with slash commands on Discord. -Phoenix
    if command == "runconsole" then
        local reconstructedString = table.concat(commandArgs, " ")
        local quotedWords = {}

        for quotedWord in reconstructedString:gmatch('"(.-)"') do
            table.insert(quotedWords, quotedWord)
        end

        local playerName = quotedWords[1]
        local consoleCommand = {}

        if #quotedWords > 1 then
            for word in quotedWords[2]:gmatch("%S+") do
                table.insert(consoleCommand, word)
            end
        end

        if #consoleCommand > 0 then
            tes3mp.LogMessage(enumerations.log.INFO, "[Discord]: Running " .. command .. ' with arguments "' .. table.concat(consoleCommand, " ") .. '"')
        else
            tes3mp.LogMessage(enumerations.log.INFO, "[Discord]: Running " .. command)
        end

        if commandHandlerData then
            local handler = commandHandlerData.handler
            handler(playerPID, playerName, consoleCommand, discordReplyChannel)
        else
            tes3mp.LogMessage(enumerations.log.WARN, "[Discord]: Unrecognized command: " .. command)
        end
    else
        local commandHandlerData = commandHandlers[command]
        if commandHandlerData then
            local handler = commandHandlerData.handler
            handler(playerPID, playerName, commandArgs, discordReplyChannel)
        else
            tes3mp.LogMessage(enumerations.log.WARN, "[Discord]: Unrecognized command: " .. command)
        end
    end
end

--- Send a response to the Discord channel.
---@param discordReplyChannel string
commands.SendResponse = function(discordReplyChannel)
    goTES3MPModules.utils.sendDiscordMessage(goTES3MP.GetServerID(), discordReplyChannel, goTES3MP.GetDefaultDiscordServer(), "**Command Executed**")
end

commands.getPlayerPID = function(str)
    if tableHelper.getCount(Players) == 0 then
        return nil
    end

    local lastPid = tes3mp.GetLastPlayerId()
    if str ~= nil then
        for playerIndex = 0, lastPid do
            if Players[playerIndex] ~= nil and Players[playerIndex]:IsLoggedIn() then
                if string.lower(Players[playerIndex].name) == string.lower(str) then
                    return playerIndex
                end
            end
        end
    end
    return nil
end


return commands