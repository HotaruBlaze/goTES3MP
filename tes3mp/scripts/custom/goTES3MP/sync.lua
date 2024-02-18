-- The goTES3MPModules variable is used to store the modules obtained from goTES3MP.GetModules().
local goTES3MPModules = nil

-- The goTES3MPSync table is used to define functions related to synchronization.
local goTES3MPSync = {}

-- The syncTimerID variable is used to store the ID of the synchronization timer.
local syncTimerID = nil

-- The syncTimer variable is used to define the duration (in seconds) between synchronization updates.
local syncTimer = 30

-- Sends a synchronization message to the server.
---@param forceResync boolean indicator whether to force a resynchronization.
goTES3MPSync.sendSync = function(forceResync)
    local serverID = goTES3MP.GetServerID()
    if serverID ~= "" then
        local maxPlayers = tostring(tes3mp.GetMaxPlayers())
        local currentPlayerCount = tostring(logicHandler.GetConnectedPlayerCount())

        -- Construct the synchronization message as a JSON object.
        local messageJson = {
            jobid = goTES3MPModules.utils.generate_uuid(),
            ServerID = serverID,
            method = "Sync",
            source = "TES3MP",
            data = {
                MaxPlayers = maxPlayers,
                CurrentPlayerCount = currentPlayerCount,
                Forced = tostring(forceResync),
                Status = "Ping"
            }
        }

        -- Encode the synchronization message into a JSON string.
        local response = goTES3MPModules.utils.isJsonValidEncode(messageJson)

        -- Send the encoded synchronization message via IrcBridge.
        if response ~= nil then
            IrcBridge.SendSystemMessage(response)
        end

        WaitingForSync = true
    end

    -- Restart the synchronization timer.
    tes3mp.RestartTimer(syncTimerID, time.seconds(syncTimer))
end

-- Callback function called when a sync message is received.
---@param serverID string - The server ID of the received sync message.
---@param receivedSyncID string - The ID of the received sync message.
goTES3MPSync.gotSync = function(serverID, receivedSyncID)
    if serverID == goTES3MP.GetServerID() then
        WaitingForSync = false
    end
end

-- Validator function registered for the "OnServerInit" event.
-- Initializes the synchronization timer and sends the initial sync message.
customEventHooks.registerValidator("OnServerInit", function()
    if syncTimerID == nil then
        -- Obtain the modules from goTES3MP.GetModules().
        goTES3MPModules = goTES3MP.GetModules()

        -- Create and start the synchronization timer.
        syncTimerID = tes3mp.CreateTimer("OnSyncUpdate", time.seconds(syncTimer))
        tes3mp.StartTimer(syncTimerID)

        -- Send the initial sync message.
        goTES3MPSync.sendSync(false)
    end
end)

-- Validator function registered for the "OnServerExit" event.
-- Stops the synchronization timer.
customEventHooks.registerValidator("OnServerExit", function()
    if syncTimerID ~= nil then
        tes3mp.StopTimer(syncTimerID)
    end
end)

-- Callback function called by the synchronization timer.
function OnSyncUpdate()
    goTES3MPSync.sendSync(false)
end

-- Return the goTES3MPSync table.
return goTES3MPSync