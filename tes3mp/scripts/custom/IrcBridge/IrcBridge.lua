-- IrcBridge.lua -*-lua-*-
-- "THE BEER-WARE LICENCE" (Revision 42):
-- <mail@michael-fitzmayer.de> wrote this file.  As long as you retain
-- this notice you can do whatever you want with this stuff. If we meet
-- some day, and you think this stuff is worth it, you can buy me a beer
-- in return.  Michael Fitzmayer

require("color")
local irc = require("irc")
local cjson = require("cjson")

local goTES3MPConfig = require("custom.goTES3MP.config")
local goTES3MPSync = require("custom.goTES3MP.sync")
local goTES3MPUtils = require("custom.goTES3MP.utils")
local goTES3MPModules = nil

local IrcBridge = {}

IrcBridge.version = "v5.0.2-goTES3MP"
IrcBridge.scriptName = "IrcBridge"
IrcBridge.debugMode = false
IrcBridge.maxMessageLength = 2048

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

                local response = goTES3MPModules["utils"].isJsonValidDecode(message)

                IrcBridge.switch(response.method) {
                    ["Sync"] = function()
                        goTES3MPSync.gotSync(response.ServerID, response.SyncID)
                    end,
                    ["Command"] = function()
                        goTES3MPModules["commands"].processCommand(response.data["TargetPlayer"],response.data["Command"],response.data["CommandArgs"], response.data["replyChannel"])
                    end,
                    ["DiscordChat"] = function()
                        IrcBridge.chatMessage(response)
                    end,
                    ["VPNCheck"] = function()
                        goTES3MPModules["vpnChecker"].kickPlayer(response.data["playerpid"], response.data["kickPlayer"])
                    end,
                    default = function()
                        print("Error: "..tableHelper.getSimplePrintableTable(response))
                        print("Unknown method (" .. response.method .. ") was received.")
                    end,
                }
            end
            lastMessage = message
        end
    )

    tes3mp.RestartTimer(IRCTimerId, time.seconds(1))
end

IrcBridge.chatMessage = function(response)
    local wherefrom = color.Default .. "[" .. response.source .. "]" .. color.Default
    local finalMessage = ""

    if response.method == "DiscordChat" then
        wherefrom = config.IRCBridge.discordColor .. "[" .. response.source .. "]" .. color.Default
    end

    if response.data["RoleColor"] ~= "" and response.data["RoleColor"] ~= "" then
        local staffRole = "#" .. response.data["RoleColor"] .. "[" .. response.data["RoleName"] .. "]" .. color.Default
        finalMessage = wherefrom .. " " .. staffRole .. " " .. response.data["User"] .. ": " .. response.data["Message"] .. "\n"
    else
        finalMessage = wherefrom .. " " .. response.data["User"] .. ": " .. response.data["Message"] .. "\n"
    end
    for pid, player in pairs(Players) do
        if Players[pid] ~= nil and Players[pid]:IsLoggedIn() then
            tes3mp.SendMessage(pid, finalMessage, false)
        end
    end
end

IrcBridge.SendSystemMessage = function(message)
    if string.len(message) > IrcBridge.maxMessageLength then
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:IRCBridge] SendSystemMessage was skipped due to message length exceeding limit.")
        return
    end
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
        goTES3MPModules = goTES3MP.GetModules()
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