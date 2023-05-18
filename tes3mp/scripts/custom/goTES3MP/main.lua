local cjson = require("cjson")
-- GoTES3MPSyncID = ""
WaitingForSync = false
local goTES3MP = {}
local TES3MPOnline = false 

-- Modules
local goTES3MPUtils = require("custom.goTES3MP.utils")
local goTES3MPSync = require("custom.goTES3MP.sync")
local goTES3MPChat = require("custom.goTES3MP.chat")
local goTES3MPVPNCheck = require("custom.goTES3MP.vpnChecker")
local goTES3MPConfig = require("custom.goTES3MP.config")


local config = goTES3MPConfig.GetConfig()

goTES3MP.GetServerID = function()
    if config.goTES3MP.serverid == "" then
        config.goTES3MP.serverid = goTES3MPUtils.randomString(16) 
        DataManager.saveData("goTES3MP", goTES3MP.config)
    end
    return config.goTES3MP.serverid
end

-- goTES3MP.GetSyncID = function()
--     if GoTES3MPSyncID == "" then
--         GoTES3MPSyncID = goTES3MPUtils.randomString(16)   
--     end
--     return GoTES3MPSyncID
-- end

goTES3MP.GetDefaultDiscordChannel = function()
    return config.goTES3MP.defaultDiscordChannel
end

goTES3MP.GetDefaultDiscordNotificationsChannel = function()
    return config.goTES3MP.defaultDiscordNotifications
end

goTES3MP.GetDefaultDiscordServer = function()
    return config.goTES3MP.defaultDiscordServer
end

customEventHooks.registerValidator(
    "OnServerInit",
    function()
        goTES3MPConfig.LoadConfig()
        goTES3MP.GetServerID()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP]: Loaded")
    end
)

customEventHooks.registerHandler("OnServerPostInit", function(eventStatus, pid)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = config.goTES3MP.serverid,
        syncid = GoTES3MPSyncID,
        data = {
            channel = config.goTES3MP.defaultDiscordNotifications,
			server = config.goTES3MP.defaultDiscordServer,
			message = "**".."[TES3MP] Server is online. :yellow_heart:".."**"
        }
    }
    if TES3MPOnline == false then
        local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
        if responce ~= nil then
            IrcBridge.SendSystemMessage(responce)
        end
        TES3MPOnline = true
    end
end)

customEventHooks.registerHandler("OnServerExit", function(eventStatus, pid)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = config.goTES3MP.serverid,
        syncid = GoTES3MPSyncID,
        data = {
            channel = config.goTES3MP.defaultDiscordNotifications,
			server = config.goTES3MP.defaultDiscordServer,
			message = "**".."[TES3MP] Server is offline. :warning:".."**"
        }
    }
    local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
    if responce ~= nil then
        IrcBridge.SendSystemMessage(responce)
    end
end)

customCommandHooks.registerCommand("forceSync", function(pid) 
    goTES3MPSync.SendSync(true)
end)

return goTES3MP