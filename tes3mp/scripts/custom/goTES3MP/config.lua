local cjson = require("cjson")
local goTES3MPConfig = {}

local config = {}
local configFile = "custom/goTES3MPConfig.json"

local defaultConfig = {
    goTES3MP = {
        serverid = "", -- Server ID
        configVersion = 1, -- Configuration version
        defaultDiscordServer = "", -- Default Discord server
        defaultDiscordChannel = "", -- Default Discord channel
        defaultDiscordNotifications = "", -- Default Discord notifications
        userModules = {}, -- User modules
    },
    IRCBridge = {
        nick = "", -- IRC nickname
        server = "", -- IRC server
        port = "6667", -- IRC port
        password = "", -- IRC password
        nspasswd = "", -- NS password
        systemchannel = "#", -- IRC system channel
        nickfilter = "", -- Nick filter
        discordColor = "#825eed", -- Discord color
        ircColor = "#5D9BEE" -- IRC color
    }
}

--- Retrieves the configuration.
---@return table - The configuration
goTES3MPConfig.GetConfig = function()
    if next(config) == nil then
        config = goTES3MPConfig.LoadConfig()
        return config
    end
    return config
end

--- Saves the configuration to the specified file.
--- @param config table - The configuration to save.
goTES3MPConfig.SaveConfig = function(config)
    if config ~= nil then
        jsonInterface.quicksave(configFile, config)
    end
end

--- Loads the configuration from the file.
--- If the file doesn't exist, a default configuration is generated.
--- If a migration is possible, the configuration is migrated.
--- @return table - The loaded configuration.
goTES3MPConfig.LoadConfig = function()
    config = jsonInterface.load(configFile)

    if config == nil then
        -- Set a default Config
        config = defaultConfig

        -- Lets check if a migration is Possible.
        local newConfig = goTES3MPConfig.MigrateFromDataManager(config)

        -- If Migration isn't possible.
        if newConfig == nil then
            -- Unable to migrate.
            goTES3MPConfig.SaveConfig(defaultConfig)
            tes3mp.LogMessage(enumerations.log.WARN, "[GoTES3MP] Migration from an old config was attempted, however failed.")
            tes3mp.LogMessage(enumerations.log.WARN, "[GoTES3MP] Default configuration has been generated at: \""..tes3mp.GetDataPath() .. "/"..configFile.."\"")
            
            tes3mp.StopServer(0)
        else
            -- Migration was Successful.
            goTES3MPConfig.SaveConfig(newConfig)
        end
    end
    return config
end

--- Migrates the configuration from the deprecated DataManager format.
--- @param config table - The current configuration.
--- @return table - The migrated configuration.
goTES3MPConfig.MigrateFromDataManager = function(config)
    -- Config file does not already exist, Lets see if we can migrate
    local dataManagerIRCConfig = jsonInterface.load("custom/__config_IrcBridge.json")
    local dataManagerGoTES3MPData = jsonInterface.load("custom/__data_goTES3MP.json")

    local newConfig = defaultConfig

    if dataManagerIRCConfig ~= nil or dataManagerGoTES3MPData ~= nil then
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:config]: Attempting to Migrate from DataManager")

        -- This was a config change before deprecating DataManager, so we need to run two different migrations.
        dataManagerGoTES3MPData = goTES3MPConfig.goTES3MPDataMigration(dataManagerGoTES3MPData)

        -- if IRCBridge config exists
        if dataManagerIRCConfig ~= nil then
            -- Migrating IRCBridge to the new single File
            for settingName, SettingValue in pairs(dataManagerIRCConfig) do
                newConfig.IRCBridge[settingName] = SettingValue
            end

            tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:config]: Migrated IRCBridge settings to new Configuration File.")
        end

        -- if goTES3MP config exists
        if dataManagerGoTES3MPData ~= nil then
            -- Migrating GoTES3MP Data to the new single File
            for settingName, SettingValue in pairs(dataManagerGoTES3MPData) do
                newConfig.goTES3MP[settingName] = SettingValue
            end
            tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:config]: Migrated GoTES3MP settings to new Configuration File.")
        end
        return newConfig
    else
        return nil
    end
end

--- Migrates the `goTES3MPData` from the deprecated DataManager format.
--- @param dataManagerGoTES3MPData table - The `goTES3MPData` from the deprecated DataManager format.
--- @return table - The migrated `goTES3MPData`.
goTES3MPConfig.goTES3MPDataMigration = function(currentConfig)
    -- Before a config version was added, we dont need to check the version currently.
    if currentConfig.configVersion == nil or currentConfig.discordchannel ~= nil or currentConfig.discordalerts ~= nil or currentConfig.discordserver ~= nil then
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:config]: Running Migration for Config N/A to Config v1.")
        local newConfig = {}

        newConfig.defaultDiscordServer = currentConfig.discordserver or currentConfig.defaultDiscordServer or ""
        newConfig.defaultDiscordChannel = currentConfig.discordchannel or currentConfig.defaultDiscordChannel or ""
        newConfig.defaultDiscordNotifications = currentConfig.discordalerts or currentConfig.defaultDiscordNotifications or ""
        newConfig.configVersion = 1
        newConfig.serverid = currentConfig.serverid or ""
        return newConfig
    else
        return currentConfig
    end

end

return goTES3MPConfig