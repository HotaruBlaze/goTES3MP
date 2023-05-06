local Sync = {}
SyncTimerID = nil

-- SyncTimer: In seconds
local SyncTimer = 30

local cjson = require("cjson")

-- Ping GoTES3MP with stats
Sync.SendSync = function(forceResync)
    local ServerID = goTES3MP.GetServerID()
    if goTES3MP.GetServerID() ~= "" then
        local messageJson = {
            ServerID = ServerID,
            method = "Sync",
            source = "TES3MP",
            data = {
                MaxPlayers = tostring(tes3mp.GetMaxPlayers()),
                CurrentPlayerCount = tostring(logicHandler.GetConnectedPlayerCount()),
                Forced = tostring(forceResync),
                Status = "Ping",
            }
        }
        IrcBridge.SendSystemMessage(cjson.encode(messageJson))

        WaitingForSync = true
    end
    tes3mp.RestartTimer(SyncTimerID, time.seconds(SyncTimer))
end

Sync.GotSync = function(ServerID, recievedSyncID)
    if ServerID == goTES3MP.GetServerID() then
        WaitingForSync = false
    end
end

customEventHooks.registerValidator(
    "OnServerInit",
    function()
        if SyncTimerID == nil then
            SyncTimerID = tes3mp.CreateTimer("OnSyncUpdate", time.seconds(SyncTimer))
            tes3mp.StartTimer(SyncTimerID)
            Sync.SendSync(false)
        end
    end
)

customEventHooks.registerValidator(
    "OnServerExit",
    function()
        if SyncTimerID ~= nil then
            tes3mp.StopTimer(SyncTimerID)
        end
    end
)

function OnSyncUpdate()
    Sync.SendSync(false)
end

return Sync