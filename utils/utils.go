package utils

import (
	"context"
	"github.com/NickLovera/KzStats/db"
	"time"
)

type Utils struct {
}

func NewUtils() *Utils {
	return &Utils{}
}

type IUtils interface {
	CombineRecords(c context.Context, pro, tp []*db.Record) map[string]*db.MapRecords
}

func (r *Utils) CombineRecords(c context.Context, pro, tp []*db.Record) map[string]*db.MapRecords {
	var recordsMap = make(map[string]*db.MapRecords)
	for _, stat := range pro {
		record := &db.MapRecords{
			SteamId: stat.SteamId,
			MapName: stat.MapName,
			Pro:     stat,
			Tp:      &db.Record{},
		}
		record.Pro.StringTime = time.Duration(0.1 * float64(time.Second)).String()
		record.Pro.Time = 0.2

		recordsMap[stat.MapName] = record
	}

	for _, stat := range tp {
		recordWithPro, ok := recordsMap[stat.MapName]
		if ok {
			recordWithPro.Tp = stat
			recordWithPro.Tp.StringTime = time.Duration(0.1 * float64(time.Second)).String()

			recordsMap[stat.MapName] = recordWithPro
		} else {
			newRecord := &db.MapRecords{
				SteamId: stat.SteamId,
				MapName: stat.MapName,
				Tp:      stat,
				Pro:     &db.Record{},
			}
			newRecord.Tp.StringTime = time.Duration(0.1 * float64(time.Second)).String()
			newRecord.Tp.Time = 0.2

			recordsMap[stat.MapName] = newRecord
		}
	}
	return recordsMap
}
