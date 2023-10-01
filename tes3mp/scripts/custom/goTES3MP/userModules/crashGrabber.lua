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
        crashType, crashReason = crashGrabber.getPreviousError()
        
        if crashReason then
            tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:crashGrabber]: Previous crash was due to\n - "..crashReason)
            goTES3MPModules.utils.sendDiscordMessage(
                serverID,
                discordChannel,
                discordServer,
                crashGrabber.generateCrashMessage(crashType, crashReason)
            )
        else
            tes3mp.LogMessage(enumerations.log.INFO, "[goTES3MP:crashGrabber]: No previous error found.")
        end
    end
)

--- Generate a crash message
---@param crashReason string The reason for the crash
---@return string The formatted crash message
function crashGrabber.generateCrashMessage(crashType, crashReason)
    return string.format("### %s\n```\n%s\n```", crashType, crashReason)
end

--- Get the second newest file in the log folder
---@param LogFolder string The path to the log folder
---@return string The name of the second newest log file
crashGrabber.getSecondNewestFile = function(LogFolder)
    local newestFile = nil
    local secondNewestFile = nil
    local newestTimestamp = 0
    local secondNewestTimestamp = 0

    local command
    if package.config:sub(1,1) == '\\' then
        -- Windows
        command = 'dir /B /O-D "' .. LogFolder .. '"'
    else
        -- Unix-like systems
        command = 'ls -1t "' .. LogFolder .. '"'
    end

    local fileHandle = io.popen(command)
    local commandOutput = fileHandle:read("*a") -- Read the entire output

    if fileHandle:close() then
        for file in commandOutput:gmatch("[^\r\n]+") do
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
    else
        print("Failed to execute command: " .. command)
    end

    return secondNewestFile
end

--- Read errors from a log file
---@param file userdata The file handle of the log file
---@return table An array of captured errors
crashGrabber.readErrorsFromLog = function(file)
    local capturedErrors = {}
    local capturedLines = {}
    local numLinesToCapture = 5
    local lastLine = nil
    local pattern = "%[(%d%d%d%d%-%d%d%-%d%d %d%d:%d%d:%d%d)%] %[(ERR)%]: .-"

    for line in file:lines() do
        local timestamp, severity = line:match(pattern)
        if timestamp and severity == "ERR" then
            table.insert(capturedErrors, {severity = severity, line = line})
        end

        table.insert(capturedLines, line)
        if #capturedLines > numLinesToCapture then
            table.remove(capturedLines, 1)
        end
    end

    return capturedErrors, lastLine
end

--- Get the previous error from the log files
--- @return string|nil The error type, and the corresponding information or nil if no error is found
crashGrabber.getPreviousError = function()
    local errorLog = crashGrabber.getSecondNewestFile(LogFolder)

    local file = assert(io.open(LogFolder.."/"..errorLog, "r"))
    local capturedErrors, lastLines  = crashGrabber.readErrorsFromLog(file)
    file:close()

    local filePathsFound = {}
    local hasScriptError = false

    for _, errorData in ipairs(capturedErrors) do
        local filePath = errorData.line:match("%.%/%a+/.+")
        if filePath then
            table.insert(filePathsFound, {severity = errorData.severity, line = errorData.line, filePath = filePath})
        end

        local matchedScriptError = errorData.line:match("%[ERR%]: %[Script%]: Error state: false")
        if matchedScriptError then
            hasScriptError = true
        end
    end

    if #filePathsFound > 0 then
        local errorFilePath = filePathsFound[1].filePath
        return "Script Error", errorFilePath
    else
        if lastLines == nil then
            return "crashGrabber was unable to find or access the previous error reason!"
        end
        if not hasScriptError then
            return "Server did not crash natually, Last log is below", table.concat(lastLines, "\n")
        end
    end
end