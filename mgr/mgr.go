package mgr

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/NickLovera/KzStats/db"
	"github.com/NickLovera/KzStats/utils"
	rest "github.com/NickLovera/rest-utils-go"
	"log"
	"net/http"
)

type KzService struct {
	KzRepo db.IKzRepo
	Utils  utils.IUtils
}

func NewKzService(repo db.IKzRepo, utils utils.IUtils) *KzService {
	return &KzService{
		KzRepo: repo,
		Utils:  utils,
	}
}

type IRequestService interface {
	GetRecordsByPlayer(c context.Context, steamId string) ([]*db.MapRecords, error)
	UpdatedRecordsByPlayer(c context.Context, steamId string) ([]*db.MapRecords, error)
	GetMapHistoryByPlayer(c context.Context, steamId, mapName string) (*db.MapRecords, error)
	//GetJumpStats
	//GetOverallStatistics
	//GetMapStats
}

func (reqService *KzService) GetRecordsByPlayer(c context.Context, steamId string) ([]*db.MapRecords, error) {
	//TODO convert playerName to steamId
	log.Printf("Getting stored stats. SteamId: %s", steamId)
	records, err := reqService.KzRepo.GetRecordsBySteamId(c, "STEAM_1:0:155905089")
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (reqService *KzService) GetMapHistoryByPlayer(c context.Context, steamId, mapName string) (*db.MapRecords, error) {
	currrRecord, err := reqService.KzRepo.GetCurrentPlayerRecordByMap(c, steamId, mapName)
	if err != nil {
		return nil, err
	}

	history, err := reqService.KzRepo.GetMapHistory(c, currrRecord.ID.Hex())
	if err != nil {
		return nil, err
	}

	currrRecord.History = history
	return currrRecord, nil
}

func (reqService *KzService) UpdatedRecordsByPlayer(c context.Context, playerName string) ([]*db.MapRecords, error) {
	var (
		proRecords []*db.Record
		tpRecords  []*db.Record
	)
	//TODO add filter for playerName
	log.Printf("Getting Pro Records. PlayerName: %s\n", playerName)
	req, err := http.NewRequest("GET", "https://kztimerglobal.com/api/v2/records/top/?steamid64=76561198272075906&has_teleports=false&modes_list_string=kz_timer&tickrate=128&limit=2000&stage=0", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request. Err: %s\n", err)
	}

	proBody, err := rest.MakeRequest(c, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(proBody, &proRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal pro body. Err: %s", err)
	}

	log.Printf("Getting Tp Records. Player: %s\n", playerName)
	req, err = http.NewRequest("GET", "https://kztimerglobal.com/api/v2/records/top/?steamid64=76561198272075906&has_teleports=true&modes_list_string=kz_timer&tickrate=128&limit=2000&stage=0", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request. Err: %s\n", err)
	}

	tpBody, err := rest.MakeRequest(c, req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tpBody, &tpRecords)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tp body. Err: %s", err)
	}

	log.Printf("Combining Records. Player: %s\n", playerName)
	combinedRecords := reqService.Utils.CombineRecords(c, proRecords, tpRecords)

	var response []*db.MapRecords
	for _, record := range combinedRecords {
		log.Println(record)
		currRecord, isPb, isNew, pbErr := reqService.KzRepo.IsPbOrNew(c, record)
		if pbErr != nil {
			return nil, pbErr
		}

		if isPb {
			log.Printf("Updating Record. Player: %s MapName: %s\n", playerName, currRecord.MapName)
			//Create new record
			newRecord, crteErr := reqService.KzRepo.CreatePlayerRecord(c, record)
			if crteErr != nil {
				return nil, crteErr
			}

			//UpdateOldRecord's parentId for history
			if _, uptdErr := reqService.KzRepo.UpdateRecordByObjectId(c, currRecord, newRecord.ID.Hex()); err != nil {
				return nil, uptdErr
			}
		} else if isNew {
			//Create new map record
			log.Printf("Creating New Record. Player: %s MapName: %s\n", playerName, record.MapName)
			_, crteErr := reqService.KzRepo.CreatePlayerRecord(c, record)
			if crteErr != nil {
				return nil, crteErr
			}
		} else {
			log.Printf("Did not beat previous pb. MapName: %s. PlayerName: %s.\n", record.MapName, playerName)
			continue
		}
		response = append(response, record)
	}

	return response, nil
}
