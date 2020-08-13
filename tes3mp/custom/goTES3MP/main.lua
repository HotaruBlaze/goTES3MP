
-- This file is used to make goTES3MP work easier with TES3MP's "Interesting" options

local cjson = require("cjson")

customEventHooks.registerHandler("OnPlayerAuthentified", function(eventStatus, pid)
	local messageJson = {
		user = tes3mp.GetName(pid),
		pid = pid,
		method = "Player",
		responce = "Connected"
	}
	IrcBridge.SendSystemMessage(cjson.encode(messageJson))
end)

customEventHooks.registerValidator("OnPlayerDisconnect", function(eventStatus, pid)
	local messageJson = {
		user = tes3mp.GetName(pid),
		pid = pid,
		method = "Player",
		responce = "Disconnected"
	}
	IrcBridge.SendSystemMessage(cjson.encode(messageJson))
end)