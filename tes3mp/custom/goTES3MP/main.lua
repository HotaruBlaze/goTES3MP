local cjson = require("cjson")
GoTES3MP_DiscordChannel = ""
GoTES3MP_DiscordServer = ""
-- GoTES3MPSyncID = ""
GOTES3MPServerID = ""
WaitingForSync = false
local goTES3MP = {}

-- Modules
local goTES3MPUtils = require("custom.goTES3MP.utils")
local goTES3MPSync = require("custom.goTES3MP.sync")
local goTES3MPChat = require("custom.goTES3MP.chat")

goTES3MP.defaultConfig = {
    serverid = "",
    discordchannel = "",
    discordserver = ""
}

goTES3MP.config = DataManager.loadData("goTES3MP", goTES3MP.defaultConfig)

goTES3MP.GetServerID = function()
    if GOTES3MPServerID == "" then
        GOTES3MPServerID = goTES3MP.config.serverid
    end
    return GOTES3MPServerID
end

-- goTES3MP.GetSyncID = function()
--     if GoTES3MPSyncID == "" then
--         GoTES3MPSyncID = goTES3MPUtils.randomString(16)   
--     end
--     return GoTES3MPSyncID
-- end

goTES3MP.GetDiscordChannel = function()
    return GoTES3MP_DiscordChannel
end
goTES3MP.GetDiscordServer = function()
    return GoTES3MP_DiscordServer
end


customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        goTES3MP.GetServerID()
        -- goTES3MP.GetSyncID()
        GoTES3MP_DiscordServer = goTES3MP.config.discordserver
        GoTES3MP_DiscordChannel = goTES3MP.config.discordchannel
    end
)


customCommandHooks.registerCommand("forceSync", function(pid) 
    goTES3MPSync.SendSync(true)
end)

return goTES3MP