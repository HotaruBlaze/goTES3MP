-- IrcBridge.lua -*-lua-*-
-- "THE BEER-WARE LICENCE" (Revision 42):
-- <mail@michael-fitzmayer.de> wrote this file.  As long as you retain
-- this notice you can do whatever you want with this stuff. If we meet
-- some day, and you think this stuff is worth it, you can buy me a beer
-- in return.  Michael Fitzmayer

require("color")
local irc = require("irc")
local cjson = require("cjson")

local goTES3MPSync = require("custom.goTES3MP.sync")
local goTES3MPUtils = require("custom.goTES3MP.utils")
local goTES3MPCommands = require("custom.goTES3MP.commands")
local goTES3MPConfig = require("custom.goTES3MP.config")
local goTES3MPVPNChecker = require("custom.goTES3MP.vpnChecker")

local IrcBridge = {}

IrcBridge.version = "v5.0.0-goTES3MP"
IrcBridge.scriptName = "IrcBridge"
IrcBridge.debugMode = false

local config = goTES3MPConfig.GetConfig()

if (config.IRCBridge.nick == "" or config.IRCBridge.systemchannel == "" or config.IRCBridge.systemchannel == "#") then
    tes3mp.LogMessage(
        enumerations.log.ERROR,
        "IrcBridge has not been configured correctly."
    )
    tes3mp.StopServer(0)
end

IRCTimerId = nil

local s = irc.new {nick = config.IRCBridge.nick}
if config.IRCBridge.password ~= "" then
    s:connect(
        {
            host = config.IRCBridge.server,
            port = config.IRCBridge.port,
            password = config.IRCBridge.password,
            timeout = 120,
            secure = false
        }
    )
else
    s:connect(config.IRCBridge.server, config.IRCBridge.port)
end
local nspasswd = "identify " .. config.IRCBridge.nspasswd
local systemchannel = config.IRCBridge.systemchannel
s:sendChat("NickServ", nspasswd)
s:join(systemchannel)
local lastMessage = ""

IrcBridge.RecvMessage = function()
    s:hook(
        "OnChat",
        function(user, systemchannel, message)
            if message ~= lastMessage then
                if IrcBridge.debugMode then
                    print("IRCDebug: " .. message)
                end

                local responce = goTES3MPUtils.isJsonValidDecode(message)

                IrcBridge.switch(responce.method) {
                    ["Sync"] = function()
                        goTES3MPSync.GotSync(responce.ServerID, responce.SyncID)
                    end,
                    ["Command"] = function()
                        goTES3MPCommands.processCommand(responce.data["TargetPlayer"],responce.data["Command"],responce.data["CommandArgs"], responce.data["replyChannel"])
                    end,
                    ["DiscordChat"] = function()
                        IrcBridge.chatMessage(responce)
                    end,
                    ["VPNCheck"] = function()
                        goTES3MPVPNChecker.kickPlayer(responce.data["playerpid"], responce.data["kickPlayer"])
                    end,
                    default = function()
                        print("Error: "..tableHelper.getSimplePrintableTable(responce))
                        print("Unknown method (" .. responce.method .. ") was received.")
                    end,
                }
            end
            lastMessage = message
        end
    )

    tes3mp.RestartTimer(IRCTimerId, time.seconds(1))
end

IrcBridge.chatMessage = function(responce)
    local wherefrom = color.Default .. "[" .. responce.source .. "]" .. color.Default
    local finalMessage = ""

    if responce.method == "DiscordChat" then
        wherefrom = config.IRCBridge.discordColor .. "[" .. responce.source .. "]" .. color.Default
    end

    if responce.data["RoleColor"] ~= "" and responce.data["RoleColor"] ~= "" then
        finalMessage = wherefrom .. " " .. staffRole .. " " .. responce.data["User"] .. ": " .. responce.data["Message"] .. "\n"
    else
        finalMessage = wherefrom .. " " .. responce.data["User"] .. ": " .. responce.data["Message"] .. "\n"
    end
    for pid, player in pairs(Players) do
        if Players[pid] ~= nil and Players[pid]:IsLoggedIn() then
            tes3mp.SendMessage(pid, finalMessage, false)
        end
    end
end

IrcBridge.SendSystemMessage = function(message)
    if message ~= lastMessage then
        s:sendChat(systemchannel, message)
        lastMessage = message
    end
end

function OnIRCUpdate()
    IrcBridge.RecvMessage()
    s:think()
end

-- Yes, i swapped out the chain of if statements for a switch statement, the performance loss is minimal 
-- however it severely increases code readability
IrcBridge.switch = function(value)
    return function(cases)
        local case = cases[value] or cases.default
        if case then
            return case(value)
        else
            error(string.format("Unhandled case (%s)", value), 2)
        end
    end
end

customEventHooks.registerValidator(
    "OnServerInit",
    function()
        IRCTimerId = tes3mp.CreateTimer("OnIRCUpdate", time.seconds(1))
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:IRCBridge] ".. IrcBridge.version.. " Loaded")
        tes3mp.StartTimer(IRCTimerId)
    end
)

customEventHooks.registerValidator(
    "OnServerExit",
    function()
        tes3mp.StopTimer(IRCTimerId)
        s:shutdown()
    end
)
return IrcBridge