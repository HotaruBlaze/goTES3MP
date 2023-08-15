local cjson = require("cjson")
-- GoTES3MPSyncID = ""
WaitingForSync = false
local goTES3MP = {}
goTES3MPModules = nil
local TES3MPOnline = false 

local goTES3MPConfig = require("custom.goTES3MP.config")
local config = goTES3MPConfig.GetConfig()

local function getLuaModulesFromFolder(folderPath)
    local luaFiles = {}

    local command
    if package.config:sub(1,1) == '\\' then
        -- Windows
        command = 'dir /B "' .. folderPath .. '"'
    else
        -- Unix-like systems
        command = 'ls -1 "' .. folderPath .. '"'
    end

    local fileHandle = io.popen(command)

    if fileHandle then
        for filename in fileHandle:lines() do
            local moduleName = filename:match("(.+)%.lua$")
            if moduleName then
                table.insert(luaFiles, moduleName)
            end
        end
        fileHandle:close()
    else
        print("Failed to execute command: " .. command)
    end

    return luaFiles
end


goTES3MP.LoadModules = function()
    if goTES3MPModules ~= nil then
        return goTES3MPModules
    end

    goTES3MPModules = {}

    -- Required Modules
    local requiredModules = {
        "utils",
        "sync",
        "commands"
    }

    tes3mp.LogMessage(enumerations.log.INFO, "[GoTES3MP:Module] Loading Modules...")
    -- Load required modules
    for _, moduleName in ipairs(requiredModules) do
        if moduleName ~= "main" then
            tes3mp.LogMessage(enumerations.log.INFO, "[GoTES3MP:Module] Loading Required Module: \""..moduleName.."\"")
            goTES3MPModules[moduleName] = require("custom.goTES3MP." .. moduleName)
        end
    end

    -- Ensure userModulesConfig exists in config
    if config["goTES3MP"]["userModules"] == nil then
        config["goTES3MP"]["userModules"] = {}
        goTES3MPConfig.SaveConfig(config)
    end

    local userModulesConfig = config["goTES3MP"]["userModules"]

    tes3mp.LogMessage(enumerations.log.INFO, "[GoTES3MP:Module] Loading user Modules...")
    -- Load user-controllable modules
    for _, moduleName in ipairs(getLuaModulesFromFolder("server/scripts/custom/goTES3MP/userModules")) do
        local moduleValue = userModulesConfig[moduleName]
        if moduleValue == true then
            tes3mp.LogMessage(enumerations.log.INFO, "[GoTES3MP:Module] Loading userModule: \""..moduleName.."\"")
            goTES3MPModules[moduleName] = require("custom.goTES3MP.userModules." .. moduleName)
        end
        userModulesConfig[moduleName] = moduleValue or false
    end

    -- Write Config
    goTES3MPConfig.SaveConfig(config)
    return goTES3MPModules
end

goTES3MP.GetServerID = function()
    if config.goTES3MP.serverid == "" then
        config.goTES3MP.serverid = goTES3MPModules["utils"].randomString(16) 
        DataManager.saveData("goTES3MP", goTES3MP.config)
    end
    return config.goTES3MP.serverid
end

goTES3MP.GetModules = function()
    if goTES3MPModules == nil then
        goTES3MPModules = goTES3MP.LoadModules()
    end
    return goTES3MPModules
end

-- goTES3MP.GetSyncID = function()
--     if GoTES3MPSyncID == "" then
--         GoTES3MPSyncID = goTES3MPModules["utils"].randomString(16)   
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
        goTES3MPModules = goTES3MP.LoadModules()
        goTES3MP.GetServerID()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP]: Loaded")
    end
)

customEventHooks.registerHandler("OnServerPostInit", function(eventStatus, pid)
    if TES3MPOnline == false then
        goTES3MPModules["utils"].sendDiscordMessage(
            config.goTES3MP.serverid,
            config.goTES3MP.defaultDiscordNotifications,
            config.goTES3MP.defaultDiscordServer,
            "**".."[TES3MP] Server is online. :yellow_heart:".."**"
        )
        TES3MPOnline = true
    end
end)

customEventHooks.registerHandler("OnServerExit", function(eventStatus, pid)
    goTES3MPModules["utils"].sendDiscordMessage(
        config.goTES3MP.serverid,
        config.goTES3MP.defaultDiscordNotifications,
        config.goTES3MP.defaultDiscordServer,
        "**".."[TES3MP] Server is offline. :warning:".."**"
    )
end)

customCommandHooks.registerCommand("forceSync", function(pid) 
    goTES3MPModules["sync"].sendSync(true)
end)

return goTES3MP