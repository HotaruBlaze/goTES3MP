local goTES3MPVPNChecker = {}
local cjson = require("cjson")
local goTES3MPUtils = require("custom.goTES3MP.utils")

customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:goTES3MPVPNChecker]: " .. "Loaded")
    end
)

-- Send IP to goTES3MP 
customEventHooks.registerHandler(
    "OnPlayerConnect",
    function(eventStatus, pid)

        local IP = tes3mp.GetIP(pid)
        local messageJson = {
            method = "VPNCheck",
            source = "TES3MP",
            serverid = GOTES3MPServerID,
            syncid = GoTES3MPSyncID,
            data = {
                channel = GoTES3MP_DiscordChannel,
                server = GoTES3MP_DiscordServer,
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

return goTES3MPVPNChecker