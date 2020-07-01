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

IrcBridge.version = "v3.0.0-goTES3MP"
IrcBridge.scriptName = "IrcBridge"

IrcBridge.defaultConfig = {
    nick = "",
    server = "",
    port = "6667",
    nspasswd = "",
    channel = "#",
    nickfilter = "",
    usrColor = "#7289da"
}

IrcBridge.config = DataManager.loadConfiguration(IrcBridge.scriptName, IrcBridge.defaultConfig)

if (IrcBridge.config == IrcBridge.defaultConfig) then
    tes3mp.LogMessage(enumerations.log.WARN, "IrcBridge configuration has been generated,")
    tes3mp.StopServer(0)
end

if (IrcBridge.config.nick == "" or IrcBridge.config.server == "" or IrcBridge.config.channel == "#") then
    tes3mp.LogMessage(
        enumerations.log.ERROR,
        "IrcBridge has not been configured correctly." .. "\n" .. "nick, server and channel are required."
    )
    tes3mp.StopServer(0)
end

local nick = IrcBridge.config.nick
local server = IrcBridge.config.server
local nspasswd = IrcBridge.config.nspasswd
local channel = IrcBridge.config.channel
local nickfilter = IrcBridge.config.nickfilter
local usrColor = IrcBridge.config.usrColor
local port = IrcBridge.config.port
IRCTimerId = nil

local s = irc.new {nick = nick}
s:connect(server, port)
nspasswd = "identify " .. nspasswd
s:sendChat("NickServ", nspasswd)
s:join(channel)
local lastMessage = ""

IrcBridge.RecvMessage = function()
    s:hook(
        "OnChat",
        function(user, channel, message)
            if message ~= lastMessage and tableHelper.getCount(Players) > 0 then
                local responce = cjson.decode(message)
                print("TES3MP: ", message)
                for pid, player in pairs(Players) do
                    if Players[pid] ~= nil and Players[pid]:IsLoggedIn() then
                        local wherefrom = usrColor .. "[Discord]" .. color.Default
                        if responce.role ~= "" and responce.role_color ~= "" then
                            local staffRole = "#"..responce.role_color .. "[" .. responce.role .. "]" .. color.Default
                            tes3mp.SendMessage(
                                pid,
                                wherefrom .." ".. staffRole .." "..responce.user .. ": " .. responce.responce .. "\n",
                                true
                            )
						else 
							tes3mp.SendMessage(
                                pid,
                                wherefrom  .." "..responce.user .. ": " .. responce.responce .. "\n",
                                true
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
    s:sendChat(channel, message)
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
        responce = cjson.encode(messageJson)
        IrcBridge.SendMessage(responce)
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
