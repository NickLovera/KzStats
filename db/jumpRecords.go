package db

import (
	"context"
	"fmt"
	mongo "github.com/NickLovera/mongo-utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JumpRepo struct {
	MongoClient mongo.IMongoClient
	//MongoUtils mongo.IMongoUtils
	CollName string
}

func NewJumpRepo(collName string, mongoClient mongo.IMongoClient) *JumpRepo {
	return &JumpRepo{
		MongoClient: mongoClient,
		CollName:    collName,
	}
}

type IJumpRepo interface {
	CreateJumpRecord(c context.Context, record *JumpRecord) (*JumpRecord, error)
	GetRecordBySteamId(c context.Context, bound bool, steamId string) (*JumpRecord, error)
	UpdateRecordByObjectId(c context.Context, oldRecord *JumpRecord, newId string) (bool, error)
	IsPbOrNew(c context.Context, bound bool, record *JumpRecord) (*JumpRecord, bool, bool, error)
	GetJumpHistory(c context.Context, bound bool, parentId string) ([]*JumpRecord, error)
}

func (repo *JumpRepo) CreateJumpRecord(c context.Context, record *JumpRecord) (*JumpRecord, error) {
	record.ID = primitive.NewObjectID()

	_, err := repo.MongoClient.InsertOne(c, repo.CollName, record)
	if err != nil {
		return nil, fmt.Errorf("faield to create player record. Err: %s", err)
	}

	return record, nil
}

func (repo *JumpRepo) GetRecordBySteamId(c context.Context, bound bool, steamId string) (*JumpRecord, error) {
	var record *JumpRecord
	err := repo.MongoClient.FindOne(c, &record, repo.CollName, bson.M{
		"is_crouch_bind":    bound,
		"steam_id":          steamId,
		"parent_history_id": "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find player jump records by steamId. Err: %s", err)
	}

	return record, nil
}

func (repo *JumpRepo) GetCurrentJumpRecordBySteamId(c context.Context, bound bool, steamId string) (*JumpRecord, error) {
	var record *JumpRecord
	err := repo.MongoClient.FindOne(c, &record, repo.CollName, bson.M{
		"is_crouch_bind":    bound,
		"steam_id":          steamId,
		"parent_history_id": "",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find players jump record  by map. Err: %s", err)
	}

	return record, nil
}

func (repo *JumpRepo) IsPbOrNew(c context.Context, bound bool, record *JumpRecord) (*JumpRecord, bool, bool, error) {
	var current *JumpRecord
	err := repo.MongoClient.FindOne(c, &current, repo.CollName, bson.M{
		"is_crouch_bind":    bound,
		"steam_id":          record.SteamId,
		"parent_history_id": "",
	})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, false, true, nil
		}
		return nil, false, false, fmt.Errorf("failed to get current jump pb. Err: %s", err)
	}
	if current != nil && record != nil {
		if record.Distance > current.Distance {
			return current, true, false, nil
		}
	}
	return nil, false, false, nil
}

func (repo *JumpRepo) GetJumpHistory(c context.Context, bound bool, parentId string) ([]*JumpRecord, error) {
	var records []*JumpRecord
	err := repo.MongoClient.FindMany(c, &records, repo.CollName, bson.M{
		"is_crouch_bind":    bound,
		"parent_history_id": parentId,
	}, options.Find().SetSort(bson.M{"_id": -1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find players jump records by steamId. Err: %s", err)
	}
	return records, nil
}

func (repo *JumpRepo) UpdateRecordByObjectId(c context.Context, oldRecord *JumpRecord, newId string) (bool, error) {
	oldRecord.ParentHistoryId = newId
	_, err := repo.MongoClient.FindOneAndReplace(c, repo.CollName, bson.M{"_id": oldRecord.ID}, oldRecord)
	if err != nil {
		return false, fmt.Errorf("failed to update jump record: %s", err)
	}
	return true, nil
}
