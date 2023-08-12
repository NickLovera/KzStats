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

type JumpService struct {
	JumpRepo db.IJumpRepo
	Utils    utils.IUtils
}

func NewJumpService(repo db.IJumpRepo, utils utils.IUtils) *JumpService {
	return &JumpService{
		JumpRepo: repo,
		Utils:    utils,
	}
}

type IJumpService interface {
	GetRecordByPlayer(c context.Context, steamId string) (*db.JumpRecord, error)
	UpdatedRecordsByPlayer(c context.Context, steamId string) ([]*db.JumpRecord, error)
	GetJumpHistoryByPlayer(c context.Context, steamId string) (*db.JumpRecord, error)
	//GetJumpStats
	//GetOverallStatistics
	//GetMapStats
}

func (jumpService *JumpService) GetRecordByPlayer(c context.Context, bound bool, steamId string) (*db.JumpRecord, error) {
	//TODO convert playerName to steamId
	log.Printf("Getting stored jumped record. SteamId: %s", steamId)
	record, err := jumpService.JumpRepo.GetRecordBySteamId(c, bound, "STEAM_1:0:155905089")
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (jumpService *JumpService) GetJumpHistoryByPlayer(c context.Context, bound bool, steamId string) (*db.JumpRecord, error) {
	currrRecord, err := jumpService.JumpRepo.GetRecordBySteamId(c, bound, steamId)
	if err != nil {
		return nil, err
	}

	history, err := jumpService.JumpRepo.GetJumpHistory(c, bound, currrRecord.ID.Hex())
	if err != nil {
		return nil, err
	}

	currrRecord.History = history
	return currrRecord, nil
}

func (jumpService *JumpService) UpdatedRecordsByPlayer(c context.Context, bound bool, playerName string) ([]*db.JumpRecord, error) {
	var (
		proRecords []*db.Record
		tpRecords  []*db.Record
	)
	//TODO add filter for playerName
	log.Printf("Getting non bind Jump Records. PlayerName: %s\n", playerName)
	req, err := http.NewRequest("GET", "https://kztimerglobal.com/api/v2/jumpstats/?jumptype_list=1&is_crouch_bind=false&is_crouch_boost=false&limit=30&steam_id=STEAM_1:0:155905089", nil)
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

	log.Printf("Getting bind jump Records. Player: %s\n", playerName)
	req, err = http.NewRequest("GET", "https://kztimerglobal.com/api/v2/jumpstats/?jumptype_list=1&is_crouch_bind=true&limit=30&steam_id=STEAM_1:0:155905089", nil)
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
	//combinedRecords := jumpService.Utils.CombineRecords(c, proRecords, tpRecords)
	//
	//var response []*db.JumpRecord
	//for _, record := range combinedRecords {
	//	currRecord, isPb, isNew, pbErr := jumpService.JumpRepo.IsPbOrNew(c, bound, record)
	//	if pbErr != nil {
	//		return nil, pbErr
	//	}
	//
	//	if isPb {
	//		log.Printf("Updating Record. Player: %s MapName: %s\n", playerName, currRecord)
	//		//Create new record
	//		newRecord, crteErr := jumpService.JumpRepo.CreatePlayerRecord(c, bound, record)
	//		if crteErr != nil {
	//			return nil, crteErr
	//		}
	//
	//		//UpdateOldRecord's parentId for history
	//		if _, uptdErr := jumpService.JumpRepo.UpdateRecordByObjectId(c, currRecord, newRecord.ID.Hex()); err != nil {
	//			return nil, uptdErr
	//		}
	//	} else if isNew {
	//		//Create new map record
	//		log.Printf("Creating New Record. Player: %s MapName: %s\n", playerName, record.MapName)
	//		//_, crteErr := jumpService.JumpRepo.CreatePlayerRecord(c, record)
	//		//if crteErr != nil {
	//		//	return nil, crteErr
	//		//}
	//	} else {
	//		//log.Printf("Did not beat previous pb. MapName: %s. PlayerName: %s.\n", record.MapName, playerName)
	//		continue
	//	}
	//	response = append(response, record)
	//}

	return nil, nil
}
