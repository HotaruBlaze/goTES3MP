local crashGrabber = {}
local cjson = require("cjson")
local goTES3MPModules = goTES3MP.GetModules()
local discordChannel = ""

local LogFolder = ".config/openmw"

customEventHooks.registerValidator(
    "OnServerPostInit",
    function()
        discordServer = goTES3MP.GetDefaultDiscordServer()
        serverID = goTES3MP.GetServerID()
        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:crashGrabber]: Loaded")

        tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:crashGrabber]: Checking if restart was due to a script error...")
        crashReason = crashGrabber.getPreviousError()
        
        if crashReason then
            tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:crashGrabber]: Previous crash was due to\n - "..crashReason)
            goTES3MPModules.utils.sendDiscordMessage(
                serverID,
                discordChannel,
                discordServer,
                crashGrabber.generateCrashMessage(crashReason)
            )
        else
            tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:crashGrabber]: No previous error found.")
        end
    end
)

--- Generate a crash message
---@param crashReason string - The reason for the crash
---@return string - The formatted crash message
function crashGrabber.generateCrashMessage(crashReason)
    return string.format("### Server Crash Detected\n```\n%s\n```", crashReason)
end

--- Get the second newest file in the log folder
---@param LogFolder string - The path to the log folder
---@return string - The name of the second newest log file
crashGrabber.getSecondNewestFile = function(LogFolder)
    local newestFile = nil
    local secondNewestFile = nil
    local newestTimestamp = 0
    local secondNewestTimestamp = 0

    if package.config:sub(1,1) == '\\' then
        -- Windows
        command = 'dir /B "' .. LogFolder .. '"'
    else
        -- Linux
        command = 'ls -lt "' .. LogFolder .. '"'
    end

    local fileHandle = io.popen(command)


    for file in io.popen(command):lines() do
        local filename = file:match("tes3mp%-server%-%d%d%d%d%-%d%d%-%d%d%-%d%d_%d%d_%d%d%.log")
        if filename then
            local timestamp = os.time { year = filename:sub(15, 18), month = filename:sub(20, 21), day = filename:sub(23, 24), hour = filename:sub(26, 27), min = filename:sub(29, 30), sec = filename:sub(32, 33) }

            if timestamp > newestTimestamp then
                secondNewestFile = newestFile
                secondNewestTimestamp = newestTimestamp
                newestFile = filename
                newestTimestamp = timestamp
            elseif timestamp > secondNewestTimestamp then
                secondNewestFile = filename
                secondNewestTimestamp = timestamp
            end
        end
    end

    return secondNewestFile
end

--- Read errors from a log file
---@param file userdata - The file handle of the log file
---@return table - An array of captured errors
crashGrabber.readErrorsFromLog = function(file)
    local capturedErrors = {}
    local pattern = "%[(%d%d%d%d%-%d%d%-%d%d %d%d:%d%d:%d%d)%] %[(ERR)%]: .-"

    for line in file:lines() do
        local timestamp, severity = line:match(pattern)
        if timestamp and severity == "ERR" then
            table.insert(capturedErrors, {severity = severity, line = line})
        end
    end

    return capturedErrors
end

--- Get the previous error from the log files
---@return string|nil - The path of the file containing the previous error, or nil if no error is found
crashGrabber.getPreviousError = function()
    local errorLog = crashGrabber.getSecondNewestFile(LogFolder)

    local file = assert(io.open(LogFolder.."/"..errorLog, "r"))
    local capturedErrors = crashGrabber.readErrorsFromLog(file)
    file:close()
    
    local filePathsFound = {}
    
    for _, errorData in ipairs(capturedErrors) do
        local filePath = errorData.line:match("%.%/%a+/.+")
        if filePath then
            table.insert(filePathsFound, {severity = errorData.severity, line = errorData.line, filePath = filePath})
        end
    end
    
    if #filePathsFound > 0 then
        local errorFilePath = filePathsFound[1].filePath
        return errorFilePath
    else
        return nil
    end
end

return crashGrabber