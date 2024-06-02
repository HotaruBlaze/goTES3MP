package main

import (
	"bytes"
	"encoding/json"
	"strconv"

	protocols "github.com/hotarublaze/gotes3mp/src/protocols"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ServerID string

func serverSync(id string, res *protocols.BaseResponse) {
	// We dont have any server saved, lets attempt to register the server.
	if viper.GetViper().GetString("tes3mp.serverid") == "" {
		if id != "" {
			viper.GetViper().Set("tes3mp.serverid", id)
		}
	}
	if viper.GetViper().GetString("tes3mp.serverid") != res.ServerId {
		if viper.GetViper().GetBool("debug") {
			log.Warnln("[DEBUG]:",
				"Ignoring'"+res.ServerId+"'",
				",Configured to use serverID",
				viper.GetViper().GetString("tes3mp.serverid"),
			)
		}
	}
	// var syncRes syncresponse
	if len(ServerID) == 0 {
		ServerID = res.Data["server_id"]
	}
	if ServerID == res.Data["server_id"] && res.Data["Status"] == "Ping" {
		if res.Data["CurrentPlayerCount"] != "" && res.Data["MaxPlayers"] != "" {
			var err error
			CurrentPlayers, err = strconv.Atoi(res.Data["CurrentPlayerCount"])
			checkError("CurrentPlayersSync", err)

			MaxPlayers, err = strconv.Atoi(res.Data["MaxPlayers"])
			checkError("MaxPlayersSync", err)
		}

		var pongresponse protocols.ServerSync

		pongresponse.ServerId = res.ServerId
		pongresponse.Status = "Pong"
		pongresponse.Method = "Sync"
		// pongresponse.SyncID = ServerSyncID

		jsonresponse, err := json.Marshal(&pongresponse)
		checkError("pongresponse", err)

		pongresponseMsg := bytes.NewBuffer(jsonresponse).String()
		IRCSendMessage(systemchannel, pongresponseMsg)
	}
}
