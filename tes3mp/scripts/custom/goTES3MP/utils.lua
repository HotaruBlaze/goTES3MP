local goTES3MPUtils = {}
local cjson = require("cjson")

local charset = {}  do -- [0-9a-zA-Z]
    for c = 48, 57  do table.insert(charset, string.char(c)) end
    for c = 65, 90  do table.insert(charset, string.char(c)) end
    for c = 97, 122 do table.insert(charset, string.char(c)) end
end

goTES3MPUtils.randomString = function(length)
    if not length or length <= 0 then return '' end
    math.randomseed(os.clock()^5)
    return goTES3MPUtils.randomString(length - 1) .. charset[math.random(1, #charset)]
end

goTES3MPUtils.isJsonValidEncode = function(json_table)
    local success, result = pcall(cjson.encode, json_table);
    if success then
        return result
    else
        return nil
    end  
end

goTES3MPUtils.isJsonValidDecode = function(json_str)
    local success, result = pcall(cjson.decode, json_str);
    if success then
        return result
    else
        return nil
    end  
end

goTES3MPUtils.sendDiscordMessage = function(ServerID, channel, server, message)
    local messageJson = {
        method = "rawDiscord",
        source = "TES3MP",
        serverid = ServerID,
        syncid = GoTES3MPSyncID,
        data = {
            channel = channel,
			server = server,
			message =message
        }
    }
    local response = goTES3MPUtils.isJsonValidEncode(messageJson)

    if response ~= nil then
        IrcBridge.SendSystemMessage(response)
    else
        tes3mp.LogMessage(enumerations.log.WARN, "[goTES3MPUtils:sendDiscordMessage] failed to send message to discord.s")
    end
end

goTES3MPUtils.alphanumsort = function(o)
	function padnum(d)
		return ("%03d%s"):format(#d, d)
	end
	table.sort(o, function(a, b)
		return tostring(a):gsub("%d+", padnum) < tostring(b):gsub("%d+", padnum)
	end)
	return o
end

return goTES3MPUtils