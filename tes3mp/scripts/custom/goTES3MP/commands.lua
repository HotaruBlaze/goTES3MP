local commands = {}
local goTES3MPModules = goTES3MP.GetModules()

-- Command handlers with descriptions
local commandHandlers = {
    ["kickplayer"] = {
        description = "Kicks the specified player from the tes3mp server.",
        handler = function(player, discordReplyChannel)
            targetPid = commands.getPlayerPID(player)
            if targetPid ~= nil then
                tes3mp.SendMessage(targetPid, color.Red .. "[SYSTEM] " .. "You have been kicked by an Administrator.", false)
                tes3mp.Kick(targetPid)
                commands.SendResponse(discordReplyChannel)
            end
        end
    },
    ["runconsole"] = {
        description = "Run a console command on a specific Player.",
        handler = function(player, commandArgs, discordReplyChannel)
            targetPid = commands.getPlayerPID(player)
            if targetPid ~= nil then
                logicHandler.RunConsoleCommandOnPlayer(targetPid, commandArgs)
                commands.SendResponse(discordReplyChannel)
            end
        end
    },
    ["resetkills"] = {
        description = "Reset player kills.",
        handler = function(player, commandArgs, discordReplyChannel)
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
        handler = function(player, commandArgs, discordReplyChannel)
            goTES3MPModules["getPlayers"].getPlayers(discordReplyChannel)
        end
    },
    ["getJournal"] = {
        description = "Get a player's Journal Entry",
        handler = function(player, commandArgs, discordReplyChannel)
            goTES3MPModules["getJournal"].GetJournalEntries(player, commandArgs, discordReplyChannel)
        end
    },
}
-- Define the !help command handler separately
local helpHandler = function(player, commandArgs, discordReplyChannel)
    local commandList = "Available commands:\n"
    for cmd, data in pairs(commandHandlers) do
        commandList = commandList .. "!" .. cmd .. ": " .. data.description .. "\n"
    end
    commandList = "```\n" .. commandList .. "```"
    
    goTES3MPModules["utils"].sendDiscordMessage(goTES3MP.GetServerID(), discordReplyChannel, goTES3MP.GetDefaultDiscordServer(), commandList)
end

-- Add the help handler to the commandHandlers table
commandHandlers["help"] = {
    description = "Display a list of available commands.",
    handler = helpHandler
}

commands.processCommand = function(player, command, commandArgs, discordReplyChannel)
    if command == nil or command == "" then
        tes3mp.LogMessage(enumerations.log.WARN, "[Discord]: processCommand triggered with blank command.")
        return
    end

    if player ~= nil then
        if string.byte(player:sub(1, 1)) == 34 then
            player = player:sub(2, string.len(player) - 1)
        end

        if commandArgs ~= nil then
            tes3mp.LogMessage(enumerations.log.INFO, "[Discord]: Running " .. command .. ' on player "' .. player .. '" with arguments "' .. commandArgs .. '"')
        else
            tes3mp.LogMessage(enumerations.log.INFO, "[Discord]: Running " .. command .. ' on player "' .. player .. '"')
        end
    else
        tes3mp.LogMessage(enumerations.log.INFO, "[Discord]: Running " .. command)
    end

    local commandHandlerData = commandHandlers[command]

    if commandHandlerData then
        local handler = commandHandlerData.handler
        handler(player, commandArgs, discordReplyChannel)
    else
        tes3mp.LogMessage(enumerations.log.WARN, "[Discord]: Unrecognized command: " .. command)
    end
end

commands.SendResponse = function(discordReplyChannel)
    goTES3MPModules["utils"].sendDiscordMessage(goTES3MP.GetServerID(), discordReplyChannel, goTES3MP.GetDefaultDiscordServer(), "**Command Executed**")
end


return commands