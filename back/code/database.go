package code

import (
	"LogC/back/models"
	"encoding/json"
	"os"
)

const filename = "logs.json"

func SaveLogsToFile(log models.Log) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(log); err != nil {
		return err
	}

	return nil
}
func ReadLogsFromFile() []models.Log {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		return nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	logs := []models.Log{}
	for {
		var log models.Log
		if err := decoder.Decode(&log); err != nil {
			break
		}
		logs = append(logs, log)
	}
	return logs
}
