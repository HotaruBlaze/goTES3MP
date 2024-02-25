local commands = {}
local goTES3MPModules = goTES3MP.GetModules()

---@class CommandHandler
---@field description string
---@field handler fun(player: string, commandArgs: string[], discordReplyChannel: string)
---@type table<string, CommandHandler>
local commandHandlers = {
    ["kickplayer"] = {
        description = "Kicks the specified player from the tes3mp server.",
        handler = function(commandArgs)
            local username = commandArgs["username"]
            local targetPid = commands.getPlayerPID(username)

            if targetPid ~= nil then
                tes3mp.SendMessage(targetPid, color.Red .. "[SYSTEM] " .. "You have been kicked by an Administrator.", false)
                tes3mp.Kick(targetPid)
                commands.sendDiscordSlashResponse(username.." has been kicked", commandArgs)
            else
                commands.sendDiscordSlashResponse(username.." does not exist", commandArgs)
            end
        end,
        args = {
            {name = "username", description = "The name of the player to kick.", required = true}
        }
    },
    ["runconsole"] = {
        description = "Run a console command on a specific Player.",
        handler = function(commandArgs)
            local username = commandArgs["username"]
            local consoleCommand = commandArgs["command"]

            local targetPid = commands.getPlayerPID(username)
            if targetPid ~= nil then
                logicHandler.RunConsoleCommandOnPlayer(targetPid, consoleCommand)
                commands.sendDiscordSlashResponse("Console command has been sent to the user", commandArgs)
            else
                commands.sendDiscordSlashResponse("Player does not exist", commandArgs)
            end
        end,
        args = {
            {name = "username", description = "The name of the player.", required = true},
            {name = "command", description = "The console command to run.", required = true}
        }
    },
    ["resetkills"] = {
        description = "Reset player kills.",
        handler = function(commandArgs)
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
            commands.sendDiscordSlashResponse("All the kill counts for creatures and NPCs have been reset", commandArgs)
        end
    },
    ["players"] = {
        description = "List Players",
        handler = function(commandArgs)
            if goTES3MPModules["getPlayers"] ~= nil then
                local playerList = goTES3MPModules.getPlayers.getPlayers()
                commands.sendDiscordSlashResponse(playerList, commandArgs)
            else
                commands.sendDiscordSlashResponse("Module not found or loaded", commandArgs)
            end
        end
    },
    ["getjournal"] = {
        description = "get a player's Journal Entry",
        handler = function(commandArgs)
            local username = commandArgs["username"]
            local questid = commandArgs["questid"]
            if goTES3MPModules["getJournal"] ~= nil then
                local questList = goTES3MPModules.getJournal.GetJournalEntries(username, questid)
                commands.sendDiscordSlashResponse(questList, commandArgs)
            else
                commands.sendDiscordSlashResponse("Module not found or loaded", commandArgs)
            end
        end,
        args = {
            {name = "username", description = "The name of the player.", required = true},
            {name = "questid", description = "the id of the quest", required = true}
        }
    },
}

--- Process the command and call the appropriate handler.
---@param command string
---@param commandArgs string[]
---@param discordReplyChannel string
commands.processCommand = function(command, commandArgs, discordReplyChannel)
    local command = string.lower(command)
    local commandHandlerData = commandHandlers[command]

    if commandHandlerData then
        local handler = commandHandlerData.handler
        handler(commandArgs)
    else
        tes3mp.LogMessage(enumerations.log.WARN, "[Discord]: Unrecognized command: " .. command)
    end
end

--- Send a response to the Discord channel.
---@param discordReplyChannel string
commands.sendDiscordSlashResponse = function(responseText, commandArgs)
    local messageJson = {
        jobid = goTES3MPModules.utils.generate_uuid(),
        method = "DiscordSlashCommandResponse",
        source = "TES3MP",
        serverid = goTES3MP.GetServerID(),
        data = commandArgs
    }
    messageJson["data"]["response"] = responseText  -- Assuming `response` is the response text

    local encodedMessage = goTES3MPModules.utils.isJsonValidEncode(messageJson)
    if encodedMessage ~= nil then
        IrcBridge.SendSystemMessage(encodedMessage)
    end
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

commands.pushSlashCommands = function(pid, cmd)
    local commandData = tableHelper.shallowCopy(commandHandlers[cmd[2]])
    
    -- Remove the "handler" field if present
    commandData.handler = nil
    -- Include the "command" field for each command
    commandData.command = cmd[2]
    
    local messageJson = {
        jobid = goTES3MPModules.utils.generate_uuid(),
        method = "RegisterDiscordSlashCommand",
        source = "TES3MP",
        serverid = goTES3MP.GetServerID(),
        data = commandData
    }

    local response = goTES3MPModules.utils.isJsonValidEncode(messageJson)
    if response ~= nil then
        IrcBridge.SendSystemMessage(response)
    end
end

commands.pushAllSlashCommands = function(pid, cmd)
    for cmdName, _ in pairs(commandHandlers) do
        local commandData = tableHelper.shallowCopy(commandHandlers[cmdName])

        -- Remove the "handler" field if present
        commandData.handler = nil
        -- Include the "command" field for each command
        commandData.command = cmdName

        local messageJson = {
            jobid = goTES3MPModules.utils.generate_uuid(),
            method = "RegisterDiscordSlashCommand",
            source = "TES3MP",
            serverid = goTES3MP.GetServerID(),
            data = commandData
        }

        local response = goTES3MPModules.utils.isJsonValidEncode(messageJson)
        if response ~= nil then
            IrcBridge.SendSystemMessage(response)
        end
    end
end

customCommandHooks.registerCommand("pushSlashCommand", function(pid, cmd) 
    commands.pushSlashCommands(pid, cmd)
end)
customCommandHooks.setRankRequirement("pushSlashCommand", 3)

customCommandHooks.registerCommand("pushAllSlashCommands", function(pid, cmd) 
    commands.pushAllSlashCommands(pid, cmd)
end)
customCommandHooks.setRankRequirement("pushAllSlashCommands", 3)


return commands