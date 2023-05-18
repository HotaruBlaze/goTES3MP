package main

import (
	"bytes"
	"encoding/json"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ServerID string

// var registrationToken
type serverSyncresponse struct {
	ServerID string `json:"serverID"`
	// SyncID   string // Removed for now
	Status string `json:"status"`
	Method string `json:"method"`
}

// type syncresponse struct {
// 	Serverid string `json:"serverid"`
// 	// SyncID             string `json:"SyncID"`
// 	MaxPlayers         int  `json:"MaxPlayers"`
// 	CurrentPlayerCount int  `json:"CurrentPlayerCount"`
// 	Forced             bool `json:"Forced"`
// 	// Players            []string `json:"Players"`
// 	// Status string `json:"Status"`
// }

func serverSync(id string, res *baseresponse) {
	// We dont have any server saved, lets attempt to register the server.
	if viper.GetViper().GetString("tes3mp.serverid") == "" {
		if id != "" {
			viper.GetViper().Set("tes3mp.serverid", id)
		}
	}
	if viper.GetViper().GetString("tes3mp.serverid") != res.ServerID {
		if viper.GetViper().GetBool("debug") {
			log.Warnln("[DEBUG]:",
				"Ignoring'"+res.ServerID+"'",
				",Configured to use serverID",
				viper.GetViper().GetString("tes3mp.serverid"),
			)
		}
	}
	// var syncRes syncresponse
	if len(ServerID) == 0 {
		ServerID = res.Data["ServerID"]
	}
	if ServerID == res.Data["ServerID"] && res.Data["Status"] == "Ping" {
		if res.Data["CurrentPlayerCount"] != "" && res.Data["MaxPlayers"] != "" {
			var err error
			CurrentPlayers, err = strconv.Atoi(res.Data["CurrentPlayerCount"])
			checkError("CurrentPlayersSync", err)

			MaxPlayers, err = strconv.Atoi(res.Data["MaxPlayers"])
			checkError("MaxPlayersSync", err)
		}

		var pongresponse serverSyncresponse

		pongresponse.ServerID = res.ServerID
		pongresponse.Status = "Pong"
		pongresponse.Method = "Sync"
		// pongresponse.SyncID = ServerSyncID

		jsonresponse, err := json.Marshal(pongresponse)
		checkError("pongresponse", err)

		pongresponseMsg := bytes.NewBuffer(jsonresponse).String()
		IRCSendMessage(systemchannel, pongresponseMsg)
	}
}
