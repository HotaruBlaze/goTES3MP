local goTES3MPVPNChecker = {}
local cjson = require("cjson")
local goTES3MPModules = goTES3MP.GetModules()

local discordChannel = ""
local discordServer = ""

local vpnWhitelist = {}

goTES3MPVPNChecker.LoadConfig = function()
    vpnWhitelist = jsonInterface.load("custom/goTES3MP_VPNWhitelist.json")

    if vpnWhitelist == nil then
        vpnWhitelist = {}
        goTES3MPVPNChecker.SaveConfig(vpnWhitelist)
    end
    return vpnWhitelist
end

goTES3MPVPNChecker.SaveConfig = function(vpnWhitelist)
    if vpnWhitelist ~= nil then
        jsonInterface.quicksave("custom/goTES3MP_VPNWhitelist.json", vpnWhitelist)
    end
end

goTES3MPVPNChecker.whitelistController = function(pid, cmd)
    if cmd[2] == "add" then
        local username = string.lower(tableHelper.concatenateFromIndex(cmd, 3))
        vpnWhitelist[string.lower(username)] = true
        goTES3MPVPNChecker.SaveConfig(vpnWhitelist)
        tes3mp.SendMessage(pid, color.RebeccaPurple .."[VPN Whitelist] " .. color.Default .. "Player \""..tableHelper.concatenateFromIndex(cmd, 3).."\" was added to the whitelist\n",false)
    end

    if cmd[2] == "remove" then
        local username = string.lower(tableHelper.concatenateFromIndex(cmd, 3))
        vpnWhitelist[string.lower(username)] = false
        goTES3MPVPNChecker.SaveConfig(vpnWhitelist)
        tes3mp.SendMessage(pid, color.RebeccaPurple .."[VPN Whitelist] " .. color.Default .. "Player \""..tableHelper.concatenateFromIndex(cmd, 3).."\" was removed from the whitelist\n",false)
    end

end

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

        local response = goTES3MPModules["utils"].isJsonValidEncode(messageJson)
        if response ~= nil then
            IrcBridge.SendSystemMessage(response)
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

            goTES3MPUtils.sendDiscordMessage(
                goTES3MP.GetServerID(),
                goTES3MP.GetDefaultDiscordChannel(),
                goTES3MP.GetDefaultDiscordServer(),
                "**"..playerName.." was kicked for trying to connect with a VPN.".."**"
            )
        end
    end
end

customCommandHooks.registerCommand("whitelist", goTES3MPVPNChecker.whitelistController)
customCommandHooks.setRankRequirement("whitelist", 1)

return goTES3MPVPNChecker