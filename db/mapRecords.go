package db

import (
	"context"
	"fmt"
	mongo "github.com/NickLovera/mongo-utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Set up repo

type KzRepo struct {
	MongoClient mongo.IMongoClient
	//MongoUtils mongo.IMongoUtils
	CollName string
}

func NewKzRepo(collName string, mongoClient mongo.IMongoClient) *KzRepo {
	return &KzRepo{
		MongoClient: mongoClient,
		CollName:    collName,
	}
}

type IKzRepo interface {
	CreatePlayerRecord(c context.Context, record *MapRecords) (*MapRecords, error)

	GetRecordsBySteamId(c context.Context, steamId string) ([]*MapRecords, error)
	GetCurrentPlayerRecordByMap(c context.Context, steamId, mapName string) (*MapRecords, error)
	GetMapHistory(c context.Context, parentId string) ([]*MapRecords, error)

	UpdateRecordByObjectId(c context.Context, oldRecord *MapRecords, newId string) (bool, error)

	IsPbOrNew(c context.Context, record *MapRecords) (*MapRecords, bool, bool, error)
	//WriteLongJumpStats()
	//WriteStatistics()
}

func (repo *KzRepo) CreatePlayerRecord(c context.Context, record *MapRecords) (*MapRecords, error) {
	record.ID = primitive.NewObjectID()

	_, err := repo.MongoClient.InsertOne(c, repo.CollName, record)
	if err != nil {
		return nil, fmt.Errorf("faield to create player record. Err: %s", err)
	}

	return record, nil
}

func (repo *KzRepo) GetCurrentPlayerRecordByMap(c context.Context, steamId, mapName string) (*MapRecords, error) {
	var record *MapRecords
	err := repo.MongoClient.FindOne(c, &record, repo.CollName, bson.M{
		"map_name":          mapName,
		"steam_id":          steamId,
		"parent_history_id": "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find player records by map. Err: %s", err)
	}

	return record, nil
}

func (repo *KzRepo) IsPbOrNew(c context.Context, record *MapRecords) (*MapRecords, bool, bool, error) {
	var current *MapRecords
	err := repo.MongoClient.FindOne(c, &current, repo.CollName, bson.M{
		"map_name":          record.MapName,
		"steam_id":          record.SteamId,
		"parent_history_id": "",
	})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, false, true, nil
		}
		return nil, false, false, fmt.Errorf("failed to get current pb. Err: %s", err)
	}
	if current.Tp != nil && record.Tp != nil {
		if current.Tp.Time > record.Tp.Time {
			return current, true, false, nil
		}
	} else if current.Pro != nil && record.Pro != nil {
		if current.Pro.Time > record.Pro.Time {
			return current, true, false, nil
		}
	}
	return nil, false, false, nil
}

func (repo *KzRepo) GetRecordsBySteamId(c context.Context, steamId string) ([]*MapRecords, error) {
	var records []*MapRecords
	err := repo.MongoClient.FindMany(c, &records, repo.CollName, bson.M{
		"steam_id":          steamId,
		"parent_history_id": "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find player records by steamId. Err: %s", err)
	}

	return records, nil
}

func (repo *KzRepo) GetMapHistory(c context.Context, parentId string) ([]*MapRecords, error) {
	var records []*MapRecords
	err := repo.MongoClient.FindMany(c, &records, repo.CollName, bson.M{
		"parent_history_id": parentId,
	}, options.Find().SetSort(bson.M{"_id": -1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find player records by steamId. Err: %s", err)
	}
	return records, nil
}

func (repo *KzRepo) UpdateRecordByObjectId(c context.Context, oldRecord *MapRecords, newId string) (bool, error) {
	oldRecord.ParentHistoryId = newId
	_, err := repo.MongoClient.FindOneAndReplace(c, repo.CollName, bson.M{"_id": oldRecord.ID}, oldRecord)
	if err != nil {
		return false, fmt.Errorf("failed to update record: %s", err)
	}
	return true, nil
}

//Record functions
//GetAllMaps
//GetRecordsByMap
//GetTpRecordsByPlayer
//GetProRecordsByPlayer

//Map functions
//GetMapById
//GetMapByDifficulty
//GetAllMaps

//JumpStats function
//GetJumpStatsByPlayer
//GetAllJumpStats

//OverallStats
// Done by comparing records to Maps I think (do in mgr)
