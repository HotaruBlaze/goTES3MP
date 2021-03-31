local commands = {}
local goTES3MPUtils = require("custom.goTES3MP.utils")

commands.kickPlayer = function(player,discordReplyChannel)
    targetPid = commands.getPlayerPID(player)
    if targetPid ~= nil then
        tes3mp.SendMessage(
            targetPid,
            color.Red .. "[SYSTEM]" .. " " .. "You have been kicked by an Administrator.",
            false
        )
        tes3mp.Kick(targetPid)
        commands.SendResponce(discordReplyChannel)
    end
end

commands.runConsole = function(player, consoleCommand, discordReplyChannel)
    targetPid = commands.getPlayerPID(player)
    if targetPid ~= nil then
        logicHandler.RunConsoleCommandOnPlayer(targetPid, consoleCommand)
        commands.SendResponce(discordReplyChannel)
    end
end

commands.resetKills = function(discordReplyChannel)
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
    commands.SendResponce(discordReplyChannel)
end

commands.main = function(str, discordReplyChannel)
    tes3mp.LogMessage(enumerations.log.INFO, "Recieved "..'"'..str..'"'.." from discord.")
    cmdChunks = {}
    for substring in str:gmatch("%S+") do
        table.insert(cmdChunks, substring)
    end

    command = string.lower(cmdChunks[1])

    if command == "kickplayer" then
        print(tableHelper.concatenateFromIndex(cmdChunks, 2))
        commands.kickPlayer(tableHelper.concatenateFromIndex(cmdChunks, 2), discordReplyChannel)
    end
    if command == "runconsole" then
        print(tableHelper.concatenateFromIndex(cmdChunks, 3))
        commands.runConsole(cmdChunks[2], tableHelper.concatenateFromIndex(cmdChunks, 3), discordReplyChannel)
    end
    if command == "resetkills" then
        commands.resetKills(discordReplyChannel)
    end
    if command == "help" then
        local commandList = ""
        commandList = commandList .. "```" .. "\n"
        commandList = commandList .. "!kickplayer (PlayerName): Kicks the specified player from the tes3mp server." .. "\n"
        commandList = commandList .. "!runconsole (PlayerName) (Console Command): Run a console command on a specific Player." .. "\n"
        commandList = commandList .. "!resetkills: Reset player kills." .. "\n"
        commandList = commandList .. "```" .. "\n"

        local messageJson = {
            method = "rawDiscord",
            source = "TES3MP",
            serverid = GOTES3MPServerID,
            syncid = GoTES3MPSyncID,
            data = {
                channel = discordReplyChannel,
                server = GoTES3MP_DiscordServer,
                message = commandList
            }
        }
    
        local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
        if responce ~= nil then
            IrcBridge.SendSystemMessage(responce)
        end
    end

end

commands.SendResponce = function(discordReplyChannel)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = GOTES3MPServerID,
        syncid = GoTES3MPSyncID,
        data = {
            channel = discordReplyChannel,
            server = GoTES3MP_DiscordServer,
            message = "**Command Executed**"
        }
    }

    local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
    if responce ~= nil then
        IrcBridge.SendSystemMessage(responce)
    end
end

commands.getPlayerPID = function(str)
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