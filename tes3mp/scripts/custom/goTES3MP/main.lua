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

goTES3MP.defaultConfig = {
    serverid = "",
    configVersion = 1,
    defaultDiscordServer = "",
    defaultDiscordChannel = "",
    defaultDiscordNotifications = "",
}

goTES3MP.config = DataManager.loadData("goTES3MP", goTES3MP.defaultConfig)

goTES3MP.GetServerID = function()
    if goTES3MP.config.serverid == "" then
        goTES3MP.config.serverid = goTES3MPUtils.randomString(16) 
        DataManager.saveData("goTES3MP", goTES3MP.config)
    end
    return goTES3MP.config.serverid
end

-- goTES3MP.GetSyncID = function()
--     if GoTES3MPSyncID == "" then
--         GoTES3MPSyncID = goTES3MPUtils.randomString(16)   
--     end
--     return GoTES3MPSyncID
-- end

goTES3MP.GetDefaultDiscordChannel = function()
    return goTES3MP.config.defaultDiscordChannel
end

goTES3MP.GetDefaultDiscordNotificationsChannel = function()
    return goTES3MP.config.defaultDiscordNotifications
end

goTES3MP.GetDefaultDiscordServer = function()
    return goTES3MP.config.defaultDiscordServer
end

customEventHooks.registerValidator(
    "OnServerInit",
    function()
        goTES3MP.GetServerID()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP]: main Initialized")
    end
)

customEventHooks.registerHandler("OnServerInit", function(eventStatus, pid)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = goTES3MP.config.serverid,
        syncid = GoTES3MPSyncID,
        data = {
            channel = goTES3MP.config.defaultDiscordNotifications,
			server = goTES3MP.config.defaultDiscordServer,
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
        serverid = goTES3MP.config.serverid,
        syncid = GoTES3MPSyncID,
        data = {
            channel = goTES3MP.config.defaultDiscordNotifications,
			server = goTES3MP.config.defaultDiscordServer,
			message = "**".."[TES3MP] Server is offline. :warning:".."**"
        }
    }
    local responce = goTES3MPUtils.isJsonValidEncode(messageJson)
    if responce ~= nil then
        IrcBridge.SendSystemMessage(responce)
    end
end)

goTES3MP.DoConfigMigration = function()
    local currentConfig = goTES3MP.config

    if currentConfig.discordchannel ~= nil or currentConfig.discordalerts ~= nil or currentConfig.discordserver ~= nil then
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:main]: Running Migration for Config N/A to Config v1.")
        -- Write a new config, using the previous config, This Assumes your using a config before Versioning was added.
        local newConfig = goTES3MP.defaultConfig
        
        newConfig.defaultDiscordServer = currentConfig.discordserver
        newConfig.defaultDiscordChannel = currentConfig.discordchannel
        newConfig.defaultDiscordNotifications = currentConfig.discordalerts
        newConfig.configVersion = 1
        newConfig.serverid = currentConfig.serverid

        -- Overwrite the config.
        goTES3MP.config = nil
        goTES3MP.config = newConfig
        DataManager.saveData("goTES3MP", goTES3MP.config)
    end

end


customCommandHooks.registerCommand("forceSync", function(pid) 
    goTES3MPSync.SendSync(true)
end)

return goTES3MP