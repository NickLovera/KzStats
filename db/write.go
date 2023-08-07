package db

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type Writer struct {
}

func NewWriter() Writer {
	return Writer{}
}

type IWriter interface {
	WritePlayerRecords(fileName string, stats []*Record) error
	//WriteLongJumpStats()
	//WriteStatistics()
}

func (wr *Writer) WritePlayerRecords(fileName string, stats []*Record) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	//Write headers
	columns := []string{
		"id",
		"player_name",
		"steam_id",
		"server_id",
		"map_id",
		"mode",
		"tickrate",
		"time",
		"string_time",
		"teleports",
		"created_on",
		"updated_on",
		"updated_by",
		"record_filter_id",
		"server_name",
		"map_name",
		"points",
		"replay_id",
	}
	w.Write(columns)

	// Using WriteAll
	var data [][]string
	for _, stat := range stats {
		row := []string{
			strconv.Itoa(stat.Id),
			stat.PlayerName,
			stat.SteamId,
			strconv.Itoa(stat.ServerId),
			strconv.Itoa(stat.MapId),
			stat.Mode,
			strconv.Itoa(stat.TickRate),
			strconv.FormatFloat(stat.Time, 'E', -1, 64),
			time.Duration(stat.Time * float64(time.Second)).String(),
			strconv.Itoa(stat.Teleports),
			stat.CreatedOn,
			stat.UpdatedOn,
			strconv.Itoa(stat.UpdatedBy),
			strconv.Itoa(stat.RecordFilterId),
			stat.ServerName,
			stat.MapName,
			strconv.Itoa(stat.Points),
			strconv.Itoa(stat.ReplayId),
		}

		data = append(data, row)
	}
	err = w.WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}
