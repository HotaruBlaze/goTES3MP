-- IrcBridge.lua -*-lua-*-
-- "THE BEER-WARE LICENCE" (Revision 42):
-- <mail@michael-fitzmayer.de> wrote this file.  As long as you retain
-- this notice you can do whatever you want with this stuff. If we meet
-- some day, and you think this stuff is worth it, you can buy me a beer
-- in return.  Michael Fitzmayer

require("color")
require("irc")
local cjson = require("cjson")

local IrcBridge = {}

IrcBridge.version = "v3.0.1-goTES3MP"
IrcBridge.scriptName = "IrcBridge"

IrcBridge.defaultConfig = {
    nick = "",
    server = "",
    port = "6667",
    password = "",
    nspasswd = "",
    systemchannel = "#",
    nickfilter = "",
    discordColor = "#825eed",
    ircColor = "#5D9BEE"
}

IrcBridge.config = DataManager.loadConfiguration(IrcBridge.scriptName, IrcBridge.defaultConfig)

if (IrcBridge.config == IrcBridge.defaultConfig) then
    tes3mp.LogMessage(enumerations.log.WARN, "IrcBridge configuration has been generated,")
    tes3mp.StopServer(0)
end

if (IrcBridge.config.nick == "" or IrcBridge.config.systemchannel == "" or IrcBridge.config.systemchannel == "#") then
    tes3mp.LogMessage(
        enumerations.log.ERROR,
        "IrcBridge has not been configured correctly." .. "\n" .. "nick, server, and systemchannel are required."
    )
    tes3mp.StopServer(0)
end

local nick = IrcBridge.config.nick
local server = IrcBridge.config.server
local nspasswd = IrcBridge.config.nspasswd
local password = IrcBridge.config.password
local systemchannel = IrcBridge.config.systemchannel
local nickfilter = IrcBridge.config.nickfilter
local discordColor = IrcBridge.config.discordColor
local ircColor = IrcBridge.config.ircColor
local port = IrcBridge.config.port
IRCTimerId = nil

local s = irc.new {nick = nick}
if password ~= "" then
    s:connect({
        host = server,
        port = port,
        password = password,
        timeout = 120,
        secure = false
    })
else
    s:connect(server, port)
end
nspasswd = "identify " .. nspasswd
s:sendChat("NickServ", nspasswd)
s:join(systemchannel)
local lastMessage = ""

IrcBridge.RecvMessage = function()
    s:hook(
        "OnChat",
        function(user, systemchannel, message)
            if message ~= lastMessage and tableHelper.getCount(Players) > 0 then
                local responce = cjson.decode(message)
                for pid, player in pairs(Players) do
                    if Players[pid] ~= nil and Players[pid]:IsLoggedIn() then
                        local wherefrom = ""
                        if responce.method == "Discord" then
                            wherefrom = discordColor .. "[" .. responce.method .. "]" .. color.Default
                        elseif responce.method == "IRC" then
                            wherefrom = ircColor .. "[" .. responce.method .. "]" .. color.Default
                        else 
                            wherefrom = color.Default .. "[" .. responce.method .. "]" .. color.Default
                        end

                        if responce.role ~= "" and responce.role_color ~= "" then
                            local staffRole = "#"..responce.role_color .. "[" .. responce.role .. "]" .. color.Default
                            tes3mp.SendMessage(
                                pid,
                                wherefrom .." ".. staffRole .." "..responce.user .. ": " .. responce.responce .. "\n",
                                false
                            )
						else 
							tes3mp.SendMessage(
                                pid,
                                wherefrom  .." "..responce.user .. ": " .. responce.responce .. "\n",
                                false
                            )
						end
                    end
                end
				lastMessage = message
            end
        end
    )

    tes3mp.RestartTimer(IRCTimerId, time.seconds(1))
end

IrcBridge.SendMessage = function(message)
    s:sendChat(systemchannel, message)
end

function OnIRCUpdate()
    IrcBridge.RecvMessage()
    s:think()
end

customEventHooks.registerValidator(
    "OnPlayerSendMessage",
    function(eventStatus, pid, message)
        local messageJson = {
            user = tes3mp.GetName(pid),
            pid = pid,
            method = "Chat",
            responce = message
        }

        if message:sub(1, 1) == "/" then
		    return
        else
            responce = cjson.encode(messageJson)
            IrcBridge.SendMessage(responce)
        end
    end
)

customEventHooks.registerValidator(
    "OnServerInit",
    function()
        IRCTimerId = tes3mp.CreateTimer("OnIRCUpdate", time.seconds(1))
    end
)

customEventHooks.registerValidator(
    "OnServerInit",
    function()
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
