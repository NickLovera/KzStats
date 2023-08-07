package routes

import (
	"context"
	rest "github.com/NickLovera/rest-utils-go"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (kzs *KzStatsServer) addRoutes(router *mux.Router) {
	router.Handle("/kzstats/recordsbyplayer/{playername}", rest.Handler{H: kzs.getRecordsByPlayer}).Methods("GET")
	router.Handle("/kzstats/updatedrecordsbyplayer/{playername}", rest.Handler{H: kzs.updatedRecordsByPlayer}).Methods("GET")
	router.Handle("/kzstats/maphistorybyplayer/{playername}/{mapname}", rest.Handler{H: kzs.mapHistoryByPlayer}).Methods("GET")

}

func (kzs *KzStatsServer) getRecordsByPlayer(c context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	mapStats, err := kzs.requestService.GetRecordsByPlayer(c, params["playername"])
	if err != nil {
		return nil, rest.ErrorCode{Code: 500, Err: err}
	}

	return rest.MyResponse().Body(mapStats).TotalSize(strconv.Itoa(len(mapStats))).Build(), nil
}

func (kzs *KzStatsServer) updatedRecordsByPlayer(c context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	mapStats, err := kzs.requestService.UpdatedRecordsByPlayer(c, params["playername"])
	if err != nil {
		return nil, rest.ErrorCode{Code: 500, Err: err}
	}

	return rest.MyResponse().Body(mapStats).TotalSize(strconv.Itoa(len(mapStats))).Build(), nil
}

func (kzs *KzStatsServer) mapHistoryByPlayer(c context.Context, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	mapStats, err := kzs.requestService.GetMapHistoryByPlayer(c, params["playername"], params["mapname"])
	if err != nil {
		return nil, rest.ErrorCode{Code: 500, Err: err}
	}

	return rest.MyResponse().Body(mapStats).TotalSize(strconv.Itoa(len(mapStats.History) + 1)).Build(), nil
}
