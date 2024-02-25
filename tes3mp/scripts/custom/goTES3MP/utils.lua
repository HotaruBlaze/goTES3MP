-- Utility functions related to various operations in GoTES3MP.
local goTES3MPUtils = {}
local cjson = require("cjson")

local charset = {}  do -- [0-9a-zA-Z]
    for c = 48, 57  do table.insert(charset, string.char(c)) end
    for c = 65, 90  do table.insert(charset, string.char(c)) end
    for c = 97, 122 do table.insert(charset, string.char(c)) end
end

-- Generates a random string of the specified length using alphanumeric characters.
---@param length number The length of the random string to generate.
---@return string The generated random string.
goTES3MPUtils.randomString = function(length)
    if not length or length <= 0 then return '' end
    math.randomseed(os.clock()^5)
    return goTES3MPUtils.randomString(length - 1) .. charset[math.random(1, #charset)]
end

-- Validates and decodes a JSON string into a Lua table.
-- @param json_str The JSON string to decode into a Lua table.
-- @return The decoded Lua table if successful, or nil if an error occurs.
goTES3MPUtils.isJsonValidEncode = function(json_table)
    local success, result = pcall(cjson.encode, json_table);
    if success then
        return result
    else
        return nil
    end  
end

-- Validates and decodes a JSON string into a Lua table.
-- @param json_str The JSON string to decode into a Lua table.
-- @return The decoded Lua table if successful, or nil if an error occurs.
goTES3MPUtils.isJsonValidDecode = function(json_str)
    local success, result = pcall(cjson.decode, json_str);
    if success then
        return result
    else
        return nil
    end  
end

--- Sends a message to Discord.
---@param ServerID string The ID of the server.
---@param channel string The channel to send the message to.
---@param server string server The target discord server id.
---@param message string The message to send to Discord.
goTES3MPUtils.sendDiscordMessage = function(ServerID, channel, server, message)
    -- Max character limit Discord allows as a single message
    if string.len(message) > 2000 then
        tes3mp.LogMessage(enumerations.log.WARN, "[goTES3MPUtils:sendDiscordMessage] message is over the 2000 character limit imposed by discord. Skipping")
        return
    end
    local messageJson = {
        jobid = goTES3MPUtils.generate_uuid(),
        method = "rawDiscord",
        source = "TES3MP",
        serverid = ServerID,
        data = {
            channel = channel,
			server = server,
			message = message
        }
    }
    local response = goTES3MPUtils.isJsonValidEncode(messageJson)
    if response ~= nil then
        IrcBridge.SendSystemMessage(response)
    else
        tes3mp.LogMessage(enumerations.log.WARN, "[goTES3MPUtils:sendDiscordMessage] failed to send message to discord.")
    end
end

-- Sorts a table in an alphanumeric order.
---@param o table The table to sort.
---@return table The sorted table.
goTES3MPUtils.alphanumsort = function(o)
	function padnum(d)
		return ("%03d%s"):format(#d, d)
	end
	table.sort(o, function(a, b)
		return tostring(a):gsub("%d+", padnum) < tostring(b):gsub("%d+", padnum)
	end)
	return o
end

goTES3MPUtils.generate_uuid = function()
    -- Set a random seed based on os.clock()
    math.randomseed(os.clock()^5)
    
    return ('xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'):gsub('[xy]', function(c)
        local v = c == 'x' and math.random(0, 15) or math.random(8, 11)
        return ('%x'):format(v)
    end)
end

return goTES3MPUtils