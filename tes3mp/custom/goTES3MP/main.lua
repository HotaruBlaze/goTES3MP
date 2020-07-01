
-- This file is used to make goTES3MP work easier with TES3MP's "Interesting" options

local cjson = require("cjson")

customEventHooks.registerValidator("OnPlayerSendMessage", function(eventStatus, pid, message)
	CharName = tes3mp.GetName(pid)
	message = tostring(message)
	
	if message:sub(1, 1) == "/" then
		return
	else
        tes3mp.LogMessage(enumerations.log.INFO, "[Chat] " .. logicHandler.GetChatName(pid) .. ": " .. message)
	end
end)

customEventHooks.registerHandler("OnPlayerAuthentified", function(eventStatus, pid)
	local messageJson = {
		user = tes3mp.GetName(pid),
		pid = pid,
		method = "User",
		responce = "Connected"
	}
	tes3mp.LogMessage(enumerations.log.INFO, "[User] " .. cjson.encode(messageJson))
end)

customEventHooks.registerValidator("OnPlayerDisconnect", function(eventStatus, pid)
	local messageJson = {
		user = tes3mp.GetName(pid),
		pid = pid,
		method = "User",
		responce = "Disconnected"
	}
	tes3mp.LogMessage(enumerations.log.INFO, "[User] " .. cjson.encode(messageJson))
end)