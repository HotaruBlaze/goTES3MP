local goTES3MPVPNChecker = {}
local cjson = require("cjson")
local goTES3MPUtils = require("custom.goTES3MP.utils")

local discordChannel = ""
local discordServer = ""

local vpnWhitelist = {
    -- ["exampleUser"] = true,
}

customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        -- Get the default configs from goTES3MP
        discordServer = goTES3MP.GetDefaultDiscordServer()
        discordChannel = goTES3MP.GetDefaultDiscordChannel()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:VPNChecker] Loaded")
    end
)

-- Send IP to goTES3MP 
customEventHooks.registerHandler(
    "OnPlayerConnect",
    function(eventStatus, pid)
        local playerName = string.lower(tes3mp.GetName(pid))

        -- If player is whitelisted, dont run the check.
        if vpnWhitelist[playerName] then
            return
        end

        local IP = tes3mp.GetIP(pid)
        local messageJson = {
            method = "VPNCheck",
            source = "TES3MP",
            serverid = goTES3MP.GetServerID(),
            syncid = GoTES3MPSyncID,
            data = {
                channel = discordChannel,
                server = discordServer,
                message = IP,
                playerpid = tostring(pid)
            }
        }

        local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
        if responce ~= nil then
            IrcBridge.SendSystemMessage(responce)
        end
    end
)

goTES3MPVPNChecker.kickPlayer = function(pid, shouldKickPlayer)
    local pid = pid
    local shouldKickPlayer = shouldKickPlayer

    if shouldKickPlayer == "yes" then
        if tes3mp.GetName(pid) ~= nil then

            playerName = tes3mp.GetName(pid)
            tes3mp.SendMessage(pid, playerName .. " was kicked for trying to use a VPN.\n", true, false)
            tes3mp.Kick(pid)

            local messageJson = {
                method = "rawDiscord",
                source = "TES3MP",
                serverid = goTES3MP.GetServerID(),
                syncid = GoTES3MPSyncID,
                data = {
                    channel = goTES3MP.GetDefaultDiscordChannel(),
                    server = goTES3MP.GetDefaultDiscordServer(),
                    message = "**"..playerName.." was kicked for trying to connect with a VPN.".."**"
                }
            }
            
            local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
            if responce ~= nil then
                IrcBridge.SendSystemMessage(responce)
            end
        end
    end
end
return goTES3MPVPNChecker