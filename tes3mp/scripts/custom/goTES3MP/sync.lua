local goTES3MPModules = nil

local goTES3MPSync = {}
local syncTimerID = nil

local syncTimer = 30

goTES3MPSync.sendSync = function(forceResync)
    local serverID = goTES3MP.GetServerID()
    if serverID ~= "" then
        local maxPlayers = tostring(tes3mp.GetMaxPlayers())
        local currentPlayerCount = tostring(logicHandler.GetConnectedPlayerCount())
        
        local messageJson = {
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
        
        local response = goTES3MPModules["utils"].isJsonValidEncode(messageJson)
        if response ~= nil then
            IrcBridge.SendSystemMessage(response)
        end
        
        WaitingForSync = true
    end
    
    tes3mp.RestartTimer(syncTimerID, time.seconds(syncTimer))
end

goTES3MPSync.gotSync = function(serverID, receivedSyncID)
    if serverID == goTES3MP.GetServerID() then
        WaitingForSync = false
    end
end

customEventHooks.registerValidator("OnServerInit", function()
    if syncTimerID == nil then
        goTES3MPModules = goTES3MP.GetModules()
        syncTimerID = tes3mp.CreateTimer("OnSyncUpdate", time.seconds(syncTimer))
        tes3mp.StartTimer(syncTimerID)
        goTES3MPSync.sendSync(false)
    end
end)

customEventHooks.registerValidator("OnServerExit", function()
    if syncTimerID ~= nil then
        tes3mp.StopTimer(syncTimerID)
    end
end)

function OnSyncUpdate()
    goTES3MPSync.sendSync(false)
end

return goTES3MPSync