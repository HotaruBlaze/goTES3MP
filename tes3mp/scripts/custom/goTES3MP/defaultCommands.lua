
-- Define the addKickPlayerCommand function
goTES3MP_Command.addCommandHandler(
    "kickplayer",
    "Kicks the specified player from the tes3mp server.",
    function(commandArgs)
        if goTES3MPModules["getPlayers"] ~= nil then
            local playerList = goTES3MPModules.getPlayers.getPlayers()
            goTES3MP_Command.sendDiscordSlashResponse(playerList, commandArgs)
        else
            goTES3MP_Command.sendDiscordSlashResponse("Module not found or loaded", commandArgs)
        end
    end,
    {}
)

goTES3MP_Command.addCommandHandler(
    "runconsole",
    "Run a console command on a specific Player.",
    function(commandArgs)
        local username = commandArgs["username"]
        local consoleCommand = commandArgs["command"]
        local targetPid = goTES3MP_Command.getPlayerPID(username)
        
        if targetPid ~= nil then
            logicHandler.RunConsoleCommandOnPlayer(targetPid, consoleCommand)
            goTES3MP_Command.sendDiscordSlashResponse("Console command has been sent to the user", commandArgs)
        else
            goTES3MP_Command.sendDiscordSlashResponse("Player does not exist", commandArgs)
        end
    end,
    {
        {name = "username", description = "The name of the player.", required = true},
        {name = "command", description = "The console command to run.", required = true}
    }
)

-- Discord Account Linking Command
goTES3MP_Command.addCommandHandler(
    "linkdiscord",
    "Generate a code to link your Discord account to your in-game account.",
    function(commandArgs)
        -- Get the Discord user ID from the command arguments
        local discordUserID = commandArgs["discordUserID"]
        
        if not discordUserID then
            goTES3MP_Command.sendDiscordSlashResponse("Error: Unable to retrieve Discord user ID.", commandArgs)
            return
        end
        
        -- Check if this Discord account has any linked accounts
        local linkedNames = goTES3MPModules["discordLinking"].getPlayerNames(discordUserID)
        if linkedNames and #linkedNames > 0 then
            local message = "Your Discord account is already linked to " .. #linkedNames .. " player account(s).\n\n" ..
                           "You can still link additional accounts by using the code provided below."
        end
        
        local linkingCode = goTES3MPModules["discordLinking"].generateCode()
        
        goTES3MPModules["discordLinking"].storePendingCode(linkingCode, discordUserID)
        
        local message = "Your linking code is: **" .. linkingCode .. "**\n\n" ..
                       "Use this code in-game with the command `/linkdiscord " .. linkingCode .. "`\n" ..
                       "This code will expire in 5 minutes."
        
        goTES3MP_Command.sendDiscordSlashResponse(message, commandArgs)
    end,
    {}
)

-- Toggle Notifications Command for Discord
goTES3MP_Command.addCommandHandler(
    "togglenotifications",
    "Toggle Discord notifications on or off for your account.",
    function(commandArgs)
        -- Get the Discord user ID from the command arguments
        local discordUserID = commandArgs["discordUserID"]
        
        if not discordUserID then
            goTES3MP_Command.sendDiscordSlashResponse("Error: Unable to retrieve Discord user ID.", commandArgs)
            return
        end
        
        -- Initialize the discord linking module if not already done
        if goTES3MPModules["discordLinking"] == nil then
            goTES3MPModules["discordLinking"] = require("custom.goTES3MP.userModules.discordLinking")
            goTES3MPModules["discordLinking"].init()
        end
        
        -- Get the notification type from command arguments (default to "sales")
        local notificationType = commandArgs["type"] or "sales"
        
        -- Get current settings
        local settings = goTES3MPModules["discordLinking"].getNotificationSettings(discordUserID)
        local currentStatus = settings[notificationType]
        
        -- Toggle the setting
        local newStatus = not currentStatus
        goTES3MPModules["discordLinking"].setNotificationPreference(discordUserID, notificationType, newStatus)
        
        -- Send confirmation message
        local statusText = newStatus and "enabled" or "disabled"
        local message = notificationType .. " notifications have been " .. statusText .. " for your account."
        
        goTES3MP_Command.sendDiscordSlashResponse(message, commandArgs)
    end,
    {
        {name = "type", description = "Type of notification to toggle (sales, etc.)", required = false}
    }
)