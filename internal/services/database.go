package services

import (
	"LogC/internal/models"
	"database/sql"
	"encoding/json"
	"os"

	_ "github.com/mattn/go-sqlite3"
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
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0666)
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

const dbFilename = "logs.db"

func InitDB() error {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		date DATETIME
	);
	CREATE TABLE IF NOT EXISTS log_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		log_id INTEGER,
		type INTEGER,
		content BLOB,
		"order" INTEGER,
		FOREIGN KEY (log_id) REFERENCES logs(id)
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveLogsToDB(log models.Log) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	insertLogQuery := `INSERT INTO logs (title, date) VALUES (?, ?)`
	result, err := tx.Exec(insertLogQuery, log.Title, log.Date)
	if err != nil {
		tx.Rollback()
		return err
	}

	logID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	insertItemQuery := `INSERT INTO log_items (log_id, type, content, "order") VALUES (?, ?, ?, ?)`
	for _, item := range log.Items {
		_, err := tx.Exec(insertItemQuery, logID, int(item.Type), item.Content, item.Order)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func ReadLogsFromDB() ([]models.Log, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, title, date FROM logs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		if err := rows.Scan(&log.Id, &log.Title, &log.Date); err != nil {
			return nil, err
		}

		itemRows, err := db.Query(`SELECT id, log_id, type, content, "order" FROM log_items WHERE log_id = ?`, log.Id)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		var items []models.LogItem
		for itemRows.Next() {
			var item models.LogItem
			if err := itemRows.Scan(&item.Id, &item.LogId, &item.Type, &item.Content, &item.Order); err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		log.Items = items
		logs = append(logs, log)
	}

	return logs, nil
}
