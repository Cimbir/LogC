package store

import (
	"LogC/internal/models"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_dbFilename = "logs.db"
	_logsTable  = "logs"
	_itemsTable = "log_items"
	_dataTable  = "log_data"
)

func InitDB(dbFilename string, logsTable string, itemsTable string, dataTable string) error {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return err
	}

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		date DATETIME
	);
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		log_id INTEGER,
		type INTEGER,
		content STRING,
		"order" INTEGER,
		FOREIGN KEY (log_id) REFERENCES logs(id)
	);
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data BLOB
	);`, logsTable, itemsTable, dataTable)
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

type SQLDB[T any] struct {
	db        *sql.DB
	tableName string
}

func NewSQLDB[T any](db *sql.DB, tableName string) *SQLDB[T] {
	return &SQLDB[T]{db: db, tableName: tableName}
}

func (s *SQLDB[T]) Add(item T) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, err
	}

	// Assuming T has fields Title and Date
	var insertQuery string
	var result sql.Result
	switch v := any(item).(type) {
	case models.Log:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (title, date) VALUES (?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.Title, v.Date)
	case models.LogItem:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (log_id, type, content, "order") VALUES (?, ?, ?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.LogId, v.Type, v.Content, v.Order)
	case models.LogData:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (data) VALUES (?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.Data)
	default:
		tx.Rollback()
		return -1, fmt.Errorf("unsupported item type: %T", item)
	}
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return int(id), tx.Commit()
}

func (s *SQLDB[T]) GetAll() ([]T, error) {
	query := fmt.Sprintf(`SELECT * FROM %s`, s.tableName)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		var item T
		switch v := any(&item).(type) {
		case *models.Log:
			if err := rows.Scan(&v.Id, &v.Title, &v.Date); err != nil {
				return nil, err
			}
		case *models.LogItem:
			if err := rows.Scan(&v.Id, &v.LogId, &v.Type, &v.Content, &v.Order); err != nil {
				return nil, err
			}
		case *models.LogData:
			if err := rows.Scan(&v.Id, &v.Data); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported item type: %T", item)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *SQLDB[T]) GetByID(id int) (T, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = ?`, s.tableName)
	row := s.db.QueryRow(query, id)
	var err error

	var item T
	switch v := any(&item).(type) {
	case *models.Log:
		if err = row.Scan(&v.Id, &v.Title, &v.Date); err != nil {
			return item, err
		}
	case *models.LogItem:
		if err = row.Scan(&v.Id, &v.LogId, &v.Type, &v.Content, &v.Order); err != nil {
			return item, err
		}
	case *models.LogData:
		if err := row.Scan(&v.Id, &v.Data); err != nil {
			return item, err
		}
	default:
		return item, fmt.Errorf("unsupported item type: %T", item)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return item, fmt.Errorf("item not found")
		}
		return item, err
	}
	return item, nil
}

func (s *SQLDB[T]) Change(id int, item T) error {
	var err error
	switch v := any(item).(type) {
	case models.Log:
		updateQuery := fmt.Sprintf(`UPDATE %s SET title = ?, date = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.Title, v.Date, id)
	case models.LogItem:
		updateQuery := fmt.Sprintf(`UPDATE %s SET log_id = ?, type = ?, content = ?, "order" = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.LogId, v.Type, v.Content, v.Order, id)
	case models.LogData:
		updateQuery := fmt.Sprintf(`UPDATE %s SET data = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.Data, id)
	default:
		return fmt.Errorf("unsupported item type: %T", item)
	}
	return err
}

func (s *SQLDB[T]) GetByField(field string, value any) ([]T, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, s.tableName, field)
	rows, err := s.db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		var item T
		switch v := any(&item).(type) {
		case *models.Log:
			if err = rows.Scan(&v.Id, &v.Title, &v.Date); err != nil {
				return nil, err
			}
		case *models.LogItem:
			if err = rows.Scan(&v.Id, &v.LogId, &v.Type, &v.Content, &v.Order); err != nil {
				return nil, err
			}
		case *models.LogData:
			if err := rows.Scan(&v.Id, &v.Data); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported item type: %T", item)
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *SQLDB[T]) Remove(id int) error {
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, s.tableName)
	_, err := s.db.Exec(deleteQuery, id)
	return err
}
