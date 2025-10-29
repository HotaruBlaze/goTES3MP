-- Discord Account Linking Module for goTES3MP
local discordLinking = {}
local cjson = require("cjson")

-- Configuration
local LINKING_DATA_FILE = "custom/discordLinking.json"
local CODE_LENGTH = 6
local CODE_EXPIRY_MINUTES = 5
local CLEANUP_TIMER_MINUTES = 10  -- Check for expired codes every 10 minutes

-- Data structure
local linkingData = {
    linkedAccounts = {},  -- discordID -> {playerName1, playerName2, ...}
    nameToDiscord = {},   -- playerName -> discordID (reverse lookup)
    pendingCodes = {},    -- code -> {discordID, expiryTime}
    notificationSettings = {} -- discordID -> {sales = true/false, ...}
}

-- Timer variables
local cleanupTimerID = nil

-- Load linking data from file
function discordLinking.loadData()
    local loadedData = jsonInterface.load(LINKING_DATA_FILE)
    if loadedData ~= nil then
        -- Ensure all required fields exist in the loaded data
        linkingData = {
            linkedAccounts = loadedData.linkedAccounts or {},
            nameToDiscord = loadedData.nameToDiscord or {},
            pendingCodes = loadedData.pendingCodes or {},
            notificationSettings = loadedData.notificationSettings or {}
        }
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Loaded linking data")
    else
        discordLinking.saveData()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Created new linking data file")
    end
end

-- Save linking data to file
function discordLinking.saveData()
    jsonInterface.quicksave(LINKING_DATA_FILE, linkingData)
end

-- Generate a random alphanumeric code
function discordLinking.generateCode()
    local charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    local code = ""
    math.randomseed(os.time())
    
    for i = 1, CODE_LENGTH do
        local rand = math.random(#charset)
        code = code .. string.sub(charset, rand, rand)
    end
    
    return code
end

-- Store a pending linking code
function discordLinking.storePendingCode(code, discordID)
    local expiryTime = os.time() + (CODE_EXPIRY_MINUTES * 60)
    linkingData.pendingCodes[code] = {
        discordID = discordID,
        expiryTime = expiryTime
    }
    discordLinking.saveData()
    tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Stored pending code: " .. code .. " for Discord ID: " .. discordID)
end

-- Check if a code is valid (exists and not expired)
function discordLinking.isCodeValid(code)
    local pendingCode = linkingData.pendingCodes[code]
    if not pendingCode then
        return false, "Invalid code"
    end
    
    if os.time() > pendingCode.expiryTime then
        discordLinking.removePendingCode(code)
        return false, "Code has expired"
    end
    
    return true, pendingCode.discordID
end

-- Remove a pending code
function discordLinking.removePendingCode(code)
    linkingData.pendingCodes[code] = nil
    discordLinking.saveData()
end

-- Link a Discord account to a player name
function discordLinking.linkAccount(discordID, playerName)
    local existingDiscordID = linkingData.nameToDiscord[playerName]
    if existingDiscordID then
        return false, "Player account is already linked to another Discord account"
    end
    
    if not linkingData.linkedAccounts[discordID] then
        linkingData.linkedAccounts[discordID] = {}
    end
    
    for _, existingName in ipairs(linkingData.linkedAccounts[discordID]) do
        if existingName == playerName then
            return false, "Player account is already linked to this Discord account"
        end
    end
    
    -- Add the name to the Discord account's linked players
    table.insert(linkingData.linkedAccounts[discordID], playerName)
    
    -- Create the reverse lookup
    linkingData.nameToDiscord[playerName] = discordID
    
    discordLinking.saveData()
    tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Linked Discord ID: " .. discordID .. " to player: " .. playerName)
    return true, "Account successfully linked"
end

-- Get Discord ID for a player name
function discordLinking.getDiscordID(playerName)
    return linkingData.nameToDiscord[playerName]
end

-- Get all player names for a Discord ID
function discordLinking.getPlayerNames(discordID)
    return linkingData.linkedAccounts[discordID] or {}
end

-- Get player name for a Discord ID (returns first one for compatibility)
function discordLinking.getPlayerName(discordID)
    local names = linkingData.linkedAccounts[discordID]
    if names and #names > 0 then
        return names[1]
    end
    return nil
end

-- Clean up expired codes
function discordLinking.cleanupExpiredCodes()
    local currentTime = os.time()
    local cleanedCount = 0
    
    for code, codeData in pairs(linkingData.pendingCodes) do
        if currentTime > codeData.expiryTime then
            linkingData.pendingCodes[code] = nil
            cleanedCount = cleanedCount + 1
        end
    end
    
    if cleanedCount > 0 then
        discordLinking.saveData()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Cleaned up " .. cleanedCount .. " expired codes")
    end
    
    return cleanedCount
end

-- Get all linked accounts (for admin purposes)
function discordLinking.getAllLinkedAccounts()
    return linkingData.linkedAccounts
end

-- Unlink a specific account (for admin purposes)
function discordLinking.unlinkAccount(discordID, playerName)
    if linkingData.linkedAccounts[discordID] then
        local names = linkingData.linkedAccounts[discordID]
        for i, name in ipairs(names) do
            if name == playerName then
                table.remove(names, i)
                
                linkingData.nameToDiscord[playerName] = nil
                
                -- If no more names linked to this Discord account, remove the entry
                if #names == 0 then
                    linkingData.linkedAccounts[discordID] = nil
                end
                
                discordLinking.saveData()
                tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Unlinked Discord ID: " .. discordID .. " from player: " .. playerName)
                return true, "Account successfully unlinked"
            end
        end
        return false, "Player name not found linked to this Discord account"
    end
    return false, "Discord account is not linked"
end

-- Unlink all accounts for a Discord ID (for admin purposes)
function discordLinking.unlinkAllAccounts(discordID)
    if linkingData.linkedAccounts[discordID] then
        local names = linkingData.linkedAccounts[discordID]
        
        for _, name in ipairs(names) do
            linkingData.nameToDiscord[name] = nil
        end
        
        linkingData.linkedAccounts[discordID] = nil
        discordLinking.saveData()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Unlinked all accounts for Discord ID: " .. discordID)
        return true, "All accounts successfully unlinked"
    end
    return false, "Discord account is not linked"
end

-- Timer callback function for cleanup
function OnDiscordLinkingCleanup()
    discordLinking.cleanupExpiredCodes()
    tes3mp.RestartTimer(cleanupTimerID, time.minutes(CLEANUP_TIMER_MINUTES))
end

-- Initialize the cleanup timer
function discordLinking.initCleanupTimer()
    if cleanupTimerID == nil then
        cleanupTimerID = tes3mp.CreateTimer("OnDiscordLinkingCleanup", time.minutes(CLEANUP_TIMER_MINUTES))
        tes3mp.StartTimer(cleanupTimerID)
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Cleanup timer initialized")
    end
end

-- Stop the cleanup timer
function discordLinking.stopCleanupTimer()
    if cleanupTimerID ~= nil then
        tes3mp.StopTimer(cleanupTimerID)
        cleanupTimerID = nil
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Cleanup timer stopped")
    end
end

-- Get notification settings for a Discord ID
function discordLinking.getNotificationSettings(discordID)
    if not linkingData.notificationSettings[discordID] then
        -- Initialize with default settings
        linkingData.notificationSettings[discordID] = {
            all = true  -- Default to enabled
        }
        discordLinking.saveData()
    end
    return linkingData.notificationSettings[discordID]
end

-- Set notification preference for a Discord ID
function discordLinking.setNotificationPreference(discordID, notificationType, enabled)
    if not linkingData.notificationSettings[discordID] then
        linkingData.notificationSettings[discordID] = {}
    end
    
    linkingData.notificationSettings[discordID][notificationType] = enabled
    discordLinking.saveData()
    tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:DiscordLinking] Set " .. notificationType .. " notifications to " .. tostring(enabled) .. " for Discord ID: " .. discordID)
    return true
end

-- Check if a player has notifications enabled for a specific type
function discordLinking.hasNotificationsEnabled(playerName, notificationType)
    local discordID = linkingData.nameToDiscord[playerName]
    if not discordID then
        return false
    end
    
    local settings = discordLinking.getNotificationSettings(discordID)
    return settings[notificationType] == true
end

customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        -- Get the default configs from goTES3MP
        discordLinking.loadData()
        discordLinking.initCleanupTimer()
        OnDiscordLinkingCleanup()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:discordLinking] Loaded")
    end
)


return discordLinking