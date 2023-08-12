package db

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MapRecords struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	SteamId         string             `json:"steam_id" bson:"steam_id"`
	MapName         string             `json:"map_name" bson:"map_name"`
	Tp              *Record            `json:"tp" bson:"tp"`
	Pro             *Record            `json:"pro" bson:"pro"`
	ParentHistoryId string             `json:"parent_history_id,omitempty" bson:"parent_history_id"`
	History         []*MapRecords      `json:"history,omitempty" bson:"-"`
}

type Record struct {
	Id             int     `json:"id" bson:"id"`
	PlayerName     string  `json:"player_name" bson:"player_name"`
	SteamId        string  `json:"steam_id" bson:"-"`
	ServerId       int     `json:"server_id" bson:"server_id"`
	MapId          int     `json:"map_id" bson:"map_id"`
	Mode           string  `json:"mode" bson:"mode"`
	TickRate       int     `json:"tickrate" bson:"tickrate"`
	Time           float64 `json:"time" bson:"time"`
	StringTime     string  `json:"string_time" bson:"string_time"`
	Teleports      int     `json:"teleports" bson:"teleports"`
	CreatedOn      string  `json:"created_on" bson:"created_on"`
	UpdatedOn      string  `json:"updated_on" bson:"updated_on"`
	UpdatedBy      int     `json:"updated_by" bson:"updated_by"`
	RecordFilterId int     `json:"record_filter_id" bson:"record_filter_id"`
	ServerName     string  `json:"server_name" bson:"server_name"`
	MapName        string  `json:"map_name" bson:"-"`
	Points         int     `json:"points" bson:"points"`
	ReplayId       int     `json:"replay_id" bson:"replay_id"`
}

func (ms *MapRecords) PrettyPrint() string {
	var response string
	response += fmt.Sprintf("MapName: %s\n", "")
	if ms.Pro != nil {
		response += fmt.Sprintf("  Pro:\n%s", ms.Pro.PrettyPrint())
	}
	if ms.Tp != nil {
		response += fmt.Sprintf("  Tp:\n%s", ms.Tp.PrettyPrint())

	}
	return response
}

func (s *Record) PrettyPrint() string {
	return fmt.Sprintf("\tId: %d\n"+
		"\tPlayerName: %s\n"+
		"\tSteamId: %s\n"+
		"\tServerId: %d\n"+
		"\tMapId: %d\n"+
		"\tMode: %s\n"+
		"\tTickRate: %d\n"+
		"\tTime: %f\n"+
		"\tTeleports: %d\n"+
		"\tCreatedOn: %s\n"+
		"\tUpdatedOn: %s\n"+
		"\tUpdatedBy: %d\n"+
		"\tRecordFilterId: %d\n"+
		"\tServerName: %s\n"+
		"\tMapName: %s\n"+
		"\tPoints: %d\n"+
		"\tReplayId: %d\n",
		s.Id, s.PlayerName, s.SteamId, s.ServerId, s.MapId, s.Mode,
		s.TickRate, s.Time, s.Teleports, s.CreatedOn, s.UpdatedOn, s.UpdatedBy,
		s.RecordFilterId, s.ServerName, s.MapName, s.Points, s.ReplayId)
}

type JumpRecord struct {
	ID              primitive.ObjectID `json:"-" bson:"_id"`
	Id              int                `json:"id" bson:"id"`
	ServerId        int                `json:"server_id" bson:"serverId"`
	PlayerName      string             `json:"player_name" bson:"playerName"`
	SteamId         string             `json:"steam_id" bson:"steamId"`
	JumpType        int                `json:"jump_type" bson:"jumpType"`
	Distance        float64            `json:"distance" bson:"distance"`
	Tickrate        int                `json:"tickrate" bson:"tickrate"`
	MslCount        int                `json:"msl_count" bson:"mslCount"`
	StrafeCount     int                `json:"strafe_count"  bson:"strafeCount"`
	IsCrouchBind    int                `json:"is_crouch_bind" bson:"isCrouchBind"`
	IsForwardBind   int                `json:"is_forward_bind" bson:"isForwardBind"`
	IsCrouchBoost   int                `json:"is_crouch_boost" bson:"isCrouchBoost"`
	UpdatedById     int                `json:"updated_by_id" bson:"updatedById"`
	CreatedOn       string             `json:"created_on" bson:"createdOn"`
	UpdatedOn       string             `json:"updated_on" bson:"updatedOn"`
	ParentHistoryId string             `json:"parent_history_id,omitempty" bson:"parent_history_id"`
	History         []*JumpRecord      `json:"history,omitempty" bson:"-"`
}
