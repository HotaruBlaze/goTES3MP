-- The cjson library is required to parse JSON data.
local cjson = require("cjson")

-- WaitingForSync is a flag indicating whether the server is waiting for synchronization.
WaitingForSync = false

-- The goTES3MP table is used to define functions related to the goTES3MP module.
local goTES3MP = {}

-- The goTES3MPModules variable is used to store the modules obtained from goTES3MP.GetModules().
goTES3MPModules = nil

-- TES3MPOnline is a flag indicating whether the TES3MP server is online.
local TES3MPOnline = false

-- The goTES3MPConfig module is required to access the goTES3MP configuration.
local goTES3MPConfig = require("custom.goTES3MP.config")

-- The goTES3MPConfig module is required to access the goTES3MP configuration.
local goTES3MPUtils = require("custom.goTES3MP.utils")

local config = goTES3MPConfig.GetConfig()

-- Helper function to get a list of Lua module names from a folder.
---@param folderPath string The path of the folder to search for Lua modules.
---@return table A table containing the Lua module names found.
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

-- Function to load the modules.
--- @return table The loaded goTES3MP modules.
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

-- Function to get the server ID.
--- @return string The server ID.
goTES3MP.GetServerID = function()
    if config.goTES3MP.serverid == "" then
        config.goTES3MP.serverid = goTES3MPModules.utils.randomString(16) 
        DataManager.saveData("goTES3MP", goTES3MP.config)
    end
    return tostring(config.goTES3MP.serverid)
end

-- Function to get the loaded modules.
--- @return table The loaded goTES3MP modules.
goTES3MP.GetModules = function()
    if goTES3MPModules ~= nil then
        return goTES3MPModules
    end

    goTES3MPModules = goTES3MP.LoadModules()
    return goTES3MPModules
end

-- Function to get the default Discord channel.
--- @return string The default Discord channel.
goTES3MP.GetDefaultDiscordChannel = function()
    return tostring(config.goTES3MP.defaultDiscordChannel)
end

-- Function to get the default Discord notifications channel.
--- @return string The default Discord notifications channel.
goTES3MP.GetDefaultDiscordNotificationsChannel = function()
    return tostring(config.goTES3MP.defaultDiscordNotifications)
end

-- Function to get the default Discord server.
---@return string The default Discord server.
goTES3MP.GetDefaultDiscordServer = function()
    return tostring(config.goTES3MP.defaultDiscordServer)
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
        goTES3MPUtils.sendDiscordMessage(
            config.goTES3MP.serverid,
            config.goTES3MP.defaultDiscordNotifications,
            config.goTES3MP.defaultDiscordServer,
            "**".."[TES3MP] Server is online. :yellow_heart:".."**"
        )
        TES3MPOnline = true
    end
end)

customEventHooks.registerHandler("OnServerExit", function(eventStatus, pid)
    goTES3MPUtils.sendDiscordMessage(
        config.goTES3MP.serverid,
        config.goTES3MP.defaultDiscordNotifications,
        config.goTES3MP.defaultDiscordServer,
        "**".."[TES3MP] Server is offline. :warning:".."**"
    )
end)

customCommandHooks.registerCommand("forceSync", function(pid) 
    goTES3MPModules.sync.sendSync(true)
end)

return goTES3MP