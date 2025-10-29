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
    local commandOutput = fileHandle:read("*a") -- Read the entire output

    if fileHandle:close() then
        for filename in commandOutput:gmatch("[^\r\n]+") do
            local moduleName = filename:match("(.+)%.lua$")
            if moduleName then
                table.insert(luaFiles, moduleName)
            end
        end
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
    if config.goTES3MP.server_id == "" then
        config.goTES3MP.server_id = goTES3MPModules.utils.randomString(16) 
        goTES3MPConfig.SaveConfig(config)
    end
    return tostring(config.goTES3MP.server_id)
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

-- Discord linking command for in-game use
customCommandHooks.registerCommand("linkdiscord", function(pid, cmd)
    
    -- Get the linking code from command arguments
    local linkingCode = cmd[2]
    
    if not linkingCode then
        tes3mp.SendMessage(pid, "Usage: /linkdiscord <code>\nGet a code from Discord using the /linkdiscord command.\n", false)
        return
    end
    
    -- Validate the code
    local isValid, result = goTES3MPModules["discordLinking"].isCodeValid(linkingCode)
    
    if not isValid then
        tes3mp.SendMessage(pid, "Invalid or expired code: " .. result .. "\n", false)
        return
    end
    
    local discordID = result
    local playerName = Players[pid].data.login.name
    
    -- Attempt to link the account
    local success, message = goTES3MPModules["discordLinking"].linkAccount(discordID, playerName)
    
    if success then
        tes3mp.SendMessage(pid, "Success! Your account has been linked to Discord.\n", false)
        
        -- Remove the used code
        goTES3MPModules["discordLinking"].removePendingCode(linkingCode)
        
        -- Send confirmation to Discord (if possible)
        local confirmMessage = "Player " .. playerName .. " has successfully linked their account."
        -- This would need to be sent via the Discord system
    else
        tes3mp.SendMessage(pid, "Failed to link account: " .. message .. "\n", false)
    end
end)

-- Command to toggle Discord notifications
customCommandHooks.registerCommand("togglenotifications", function(pid, cmd)
    local playerName = Players[pid].data.login.name
    
    -- Check if the player has a linked Discord account
    local discordID = goTES3MPModules["discordLinking"].getDiscordID(playerName)
    if not discordID then
        tes3mp.SendMessage(pid, "You must link your Discord account first using /linkdiscord <code>\n", false)
        return
    end
    
    -- Get the notification type from command arguments (default to "sales")
    local notificationType = cmd[2] or "sales"
    
    -- Get current settings
    local settings = goTES3MPModules["discordLinking"].getNotificationSettings(discordID)
    local currentStatus = settings[notificationType]
    
    -- Toggle the setting
    local newStatus = not currentStatus
    goTES3MPModules["discordLinking"].setNotificationPreference(discordID, notificationType, newStatus)
    
    -- Send confirmation message
    local statusText = newStatus and "enabled" or "disabled"
    tes3mp.SendMessage(pid, notificationType .. " notifications have been " .. statusText .. ".\n", false)
end)

-- Function to send sale notification to a player
goTES3MP.sendSaleNotification = function(playerName, itemName, price, quantity)
    if goTES3MPModules == nil then
        goTES3MPModules = goTES3MP.LoadModules()
    end
    
    if goTES3MPModules["discordLinking"] == nil then
        tes3mp.LogMessage(enumerations.log.WARN, "[goTES3MP] Discord linking module not loaded")
        return false
    end
    
    -- Create the fields for the embed
    local fields = {
        {name = "Item", value = itemName, inline = true},
        {name = "Price", value = tostring(price) .. " gold", inline = true},
        {name = "Quantity", value = tostring(quantity), inline = true}
    }
    
    -- Send the notification
    local success = goTES3MPModules["discordLinking"].sendNotification(
        playerName,
        "sales",
        "💰 Item Sold!",
        "Your item has been sold successfully.",
        fields
    )
    
    if success then
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP] Sale notification sent to " .. playerName)
    else
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP] Sale notification not sent to " .. playerName .. " (notifications disabled or not linked)")
    end
    
    return success
end

return goTES3MP
