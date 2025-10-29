goTES3MP_Command = {}
local goTES3MPModules = goTES3MP.GetModules()
local loadDefaultCommands = true

---@class CommandHandler
---@field description string
---@field handler fun(player: string, commandArgs: string[], discordReplyChannel: string)
---@type table<string, CommandHandler>
local commandHandlers = {}

---@param commandName string
---@param description string
---@param handler fun(commandArgs: table)
---@param args table
goTES3MP_Command.addCommandHandler = function(commandName, description, handler, args)
    tes3mp.LogMessage(enumerations.log.INFO, "[Discord]: Adding command handler for " .. commandName)
    if commandHandlers[commandName] then
        commandHandlers[commandName] = nil
    end
    commandHandlers[commandName] = {
        command = commandName,
        description = description,
        handler = handler,
        args = args
    }
end

--- Process the command and call the appropriate handler.
---@param command string
---@param commandArgs string[]
---@param discordReplyChannel string
---@param discordUserID string (optional)
goTES3MP_Command.processCommand = function(command, commandArgs, discordReplyChannel, discordUserID)
    local command = string.lower(command)
    local commandHandlerData = commandHandlers[command]

    if commandHandlerData then
        local handler = commandHandlerData.handler
        
        -- Add discordUserID to commandArgs if provided
        if discordUserID then
            commandArgs["discordUserID"] = discordUserID
        end
        
        handler(commandArgs)
    else
        tes3mp.LogMessage(enumerations.log.WARN, "[Discord]: Unrecognized command: " .. command)
    end
end

--- Send a response to the Discord channel.
---@param discordReplyChannel string
goTES3MP_Command.sendDiscordSlashResponse = function(responseText, commandArgs)
    local messageJson = {
        job_id = goTES3MPModules.utils.generate_uuid(),
        server_id = goTES3MP.GetServerID(),
        method = "DiscordSlashCommandResponse",
        source = "TES3MP",
        data = commandArgs
    }
    messageJson["data"]["response"] = responseText  -- Assuming `response` is the response text

    local encodedMessage = goTES3MPModules.utils.isJsonValidEncode(messageJson)
    if encodedMessage ~= nil then
        IrcBridge.SendSystemMessage(encodedMessage)
    end
end


goTES3MP_Command.getPlayerPID = function(str)
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

goTES3MP_Command.pushSlashCommands = function(pid, cmd)
    local commandData = tableHelper.shallowCopy(commandHandlers[cmd[2]])
    
    -- Remove the "handler" field if present
    commandData.handler = nil
    -- Include the "command" field for each command
    commandData.command = cmd[2]
    
    local messageJson = {
        job_id = goTES3MPModules.utils.generate_uuid(),
        method = "RegisterDiscordSlashCommand",
        source = "TES3MP",
        server_id = goTES3MP.GetServerID(),
        data = commandData
    }

    local response = goTES3MPModules.utils.isJsonValidEncode(messageJson)
    if response ~= nil then
        IrcBridge.SendSystemMessage(response)
    end
end

goTES3MP_Command.pushAllSlashCommands = function(pid, cmd)
    for cmdName, _ in pairs(commandHandlers) do
        local commandData = tableHelper.shallowCopy(commandHandlers[cmdName])

        -- Remove the "handler" field if present
        commandData.handler = nil
        -- Include the "command" field for each command
        commandData.command = cmdName

        local messageJson = {
            job_id = goTES3MPModules.utils.generate_uuid(),
            method = "RegisterDiscordSlashCommand",
            source = "TES3MP",
            server_id = goTES3MP.GetServerID(),
            data = commandData
        }
        local response = goTES3MPModules.utils.isJsonValidEncode(messageJson)
        if response ~= nil then
            IrcBridge.SendSystemMessage(response)
        end
    end
end

if loadDefaultCommands then
    require("custom.goTES3MP.defaultCommands")
end

customCommandHooks.registerCommand("pushSlashCommand", function(pid, cmd) 
    goTES3MP_Command.pushSlashCommands(pid, cmd)
end)
customCommandHooks.setRankRequirement("pushSlashCommand", 3)

customCommandHooks.registerCommand("pushAllSlashCommands", function(pid, cmd) 
    goTES3MP_Command.pushAllSlashCommands(pid, cmd)
end)
customCommandHooks.setRankRequirement("pushAllSlashCommands", 3)

return goTES3MP_Command