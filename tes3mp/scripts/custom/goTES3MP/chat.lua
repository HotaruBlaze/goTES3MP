local cjson = require("cjson")
local goTES3MPUtils = require("custom.goTES3MP.utils")
local goTES3MPChat = {}

local discordServer = ""
local discordChannel = ""

local maxCharMessageLength = 1512 -- This can be set to 1512 if using my personal fork(Temporary fix) 450 Default

customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        -- Get the default configs from goTES3MP
        discordServer = goTES3MP.GetDefaultDiscordServer()
        discordChannel = goTES3MP.GetDefaultDiscordChannel()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:chat] Loaded")
    end
)

customEventHooks.registerHandler(
    "OnPlayerAuthentified",
    function(eventStatus, pid)
        goTES3MPUtils.sendDiscordMessage(
            goTES3MP.GetServerID(),
            goTES3MP.GetDefaultDiscordChannel(),
            goTES3MP.GetDefaultDiscordServer(),
            "**" .. "[TES3MP] " .. tes3mp.GetName(pid) .. " has connected" .. "**"
        )
    end
)

customEventHooks.registerValidator(
    "OnPlayerDisconnect",
    function(eventStatus, pid)
        goTES3MPUtils.sendDiscordMessage(
            goTES3MP.GetServerID(),
            goTES3MP.GetDefaultDiscordChannel(),
            goTES3MP.GetDefaultDiscordServer(),
            "**" .. "[TES3MP] " .. tes3mp.GetName(pid) .. " has disconnected" .. "**"
        )
    end)

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
                goTES3MPUtils.sendDiscordMessage(
                    goTES3MP.GetServerID(),
                    goTES3MP.GetDefaultDiscordChannel(),
                    goTES3MP.GetDefaultDiscordServer(),
                    tes3mp.GetName(pid) .. ": " .. message
                )
            end
        end
    end
)

customEventHooks.registerValidator(
    "OnPlayerDeath",
    function(eventStatus, pid)
        local playerName = Players[pid].name
        local deathReason = "committed suicide"

        if tes3mp.DoesPlayerHavePlayerKiller(pid) then
            local killerPid = tes3mp.GetPlayerKillerPid(pid)
            if pid ~= killerPid then deathReason = "was killed by player " .. Players[killerPid].name end
        else
            local killerName = tes3mp.GetPlayerKillerName(pid)
            if killerName ~= "" then deathReason = "was killed by " .. killerName end
        end

        deathReason = playerName .. " " .. deathReason

        goTES3MPUtils.sendDiscordMessage(
            goTES3MP.GetServerID(),
            goTES3MP.GetDefaultDiscordChannel(),
            goTES3MP.GetDefaultDiscordServer(),
            "***" .. deathReason .. "***"
        )
    end
)

return goTES3MPChat