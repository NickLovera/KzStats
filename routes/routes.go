package routes

import (
	rest "github.com/NickLovera/rest-utils-go"
	"github.com/gorilla/mux"
)

//TODO move to db collection
var NameToSteamId = map[string]string{
	"DingleDorf": "STEAM_1:0:155905089",
	"s":          "STEAM_1:0:155905089",
}

func (kzs *KzStatsServer) addRoutes(router *mux.Router) {
	//MapRecords
	router.Handle("/kzstats/recordsbyplayer/{playername}", rest.Handler{H: kzs.getRecordsByPlayer}).Methods("GET", "OPTIONS")
	router.Handle("/kzstats/updatedrecordsbyplayer/{playername}", rest.Handler{H: kzs.updatedRecordsByPlayer}).Methods("GET", "OPTIONS")
	router.Handle("/kzstats/maphistorybyplayer/{playername}/{mapname}", rest.Handler{H: kzs.mapHistoryByPlayer}).Methods("GET", "OPTIONS")
}
