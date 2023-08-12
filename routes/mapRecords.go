package routes

import (
	"context"
	rest "github.com/NickLovera/rest-utils-go"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (kzs *KzStatsServer) getRecordsByPlayer(c context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	mapStats, err := kzs.requestService.GetRecordsByPlayer(c, NameToSteamId[params["playername"]])
	if err != nil {
		return nil, rest.ErrorCode{Code: 500, Err: err}
	}
	return rest.MyResponse().Body(mapStats).TotalSize(strconv.Itoa(len(mapStats))).Build(), nil
}

func (kzs *KzStatsServer) updatedRecordsByPlayer(c context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	mapStats, err := kzs.requestService.UpdatedRecordsByPlayer(c, NameToSteamId[params["playername"]])
	if err != nil {
		return nil, rest.ErrorCode{Code: 500, Err: err}
	}

	return rest.MyResponse().Body(mapStats).TotalSize(strconv.Itoa(len(mapStats))).Build(), nil
}

func (kzs *KzStatsServer) mapHistoryByPlayer(c context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	mapStats, err := kzs.requestService.GetMapHistoryByPlayer(c, NameToSteamId[params["playername"]], params["mapname"])
	if err != nil {
		return nil, rest.ErrorCode{Code: 500, Err: err}
	}

	return rest.MyResponse().Body(mapStats).TotalSize(strconv.Itoa(len(mapStats.History) + 1)).Build(), nil
}
