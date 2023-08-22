local cjson = require("cjson")
local goTES3MPChat = {}
local goTES3MPModules = goTES3MP.GetModules()

local serverID = ""
local discordServer = ""
local discordChannel = ""

local maxCharMessageLength = 1512 -- This can be set to 1512 if using my personal fork(Temporary fix) 450 Default

customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        -- Get the default configs from goTES3MP
        discordServer = goTES3MP.GetDefaultDiscordServer()
        discordChannel = goTES3MP.GetDefaultDiscordChannel()
        serverID = goTES3MP.GetServerID()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:chat] Loaded")
    end
)

-- Handle player authentication event
customEventHooks.registerHandler(
    "OnPlayerAuthentified",
    function(eventStatus, pid)
        goTES3MPModules.utils.sendDiscordMessage(
            serverID,
            discordChannel,
            discordServer,
            "**" .. "[TES3MP] " .. tes3mp.GetName(pid) .. " has connected" .. "**"
        )
    end
)

-- Handle player disconnection event
customEventHooks.registerValidator(
    "OnPlayerDisconnect",
    function(eventStatus, pid)
        goTES3MPModules.utils.sendDiscordMessage(
            serverID,
            discordChannel,
            discordServer,
            "**" .. "[TES3MP] " .. tes3mp.GetName(pid) .. " has disconnected" .. "**"
        )
    end
)

-- Handle player chat message event
customEventHooks.registerValidator(
    "OnPlayerSendMessage",
    function(eventStatus, pid, message)
        if string.len(message) > maxCharMessageLength then
            tes3mp.SendMessage(pid, color.Red .."[goTES3MP] " .. color.Default .. "Warning, Message was too long and was not relayed to discord\n",false)
            tes3mp.LogMessage(enumerations.log.WARN, "Chat message for " .. '"' .. tes3mp.GetName(pid) .. '"' .. " was not sent")
        else
            if message:sub(1, 1) == "/" then
                return
            else
                goTES3MPModules.utils.sendDiscordMessage(
                    serverID,
                    discordChannel,
                    discordServer,
                    tes3mp.GetName(pid) .. ": " .. message
                )
            end
        end
    end
)

-- Handle player death event
customEventHooks.registerValidator(
    "OnPlayerDeath",
    function(eventStatus, pid)
        local playerName = Players[pid].name
        local deathReason = "committed suicide"

        if tes3mp.DoesPlayerHavePlayerKiller(pid) then
            local killerPid = tes3mp.GetPlayerKillerPid(pid)
            if pid ~= killerPid then
                deathReason = "was killed by player " .. Players[killerPid].name
            end
        else
            local killerName = tes3mp.GetPlayerKillerName(pid)
            if killerName ~= "" then
                deathReason = "was killed by " .. killerName
            end
        end

        deathReason = playerName .. " " .. deathReason

        goTES3MPModules.utils.sendDiscordMessage(
            serverID,
            discordChannel,
            discordServer,
            "***" .. deathReason .. "***"
        )
    end
)

return goTES3MPChat