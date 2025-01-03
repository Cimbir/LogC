package store

import (
	models "LogC/internal/models/store"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFilename   = "logs.db"
	LogsTable    = "logs"
	ItemsTable   = "log_items"
	DataTable    = "log_data"
	UserTable    = "users"
	CommentTable = "comments"
)

func InitDB(dbFilename, logsTable, itemsTable, dataTable, userTable, commentTable string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		return nil, err
	}

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title STRING,
		date DATETIME,
		thumbnail_id INTEGER,
		category INTEGER,
		short_desc STRING
	);
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		log_id INTEGER,
		type INTEGER,
		content TEXT,
		"order" INTEGER,
		FOREIGN KEY (log_id) REFERENCES logs(id)
	);
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data BLOB,
		desc STRING
	);
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username STRING,
		password STRING,
		is_admin BOOLEAN
	);
	CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		log_id INTEGER,
		content TEXT,
		date DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (log_id) REFERENCES logs(id)
	);`, logsTable, itemsTable, dataTable, userTable, commentTable)
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
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
		insertQuery = fmt.Sprintf(`INSERT INTO %s (title, date, thumbnail_id, category, short_desc) VALUES (?, ?, ?, ?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.Title, v.Date, v.ThumbnailId, v.Category, v.ShortDesc)
	case models.LogItem:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (log_id, type, content, "order") VALUES (?, ?, ?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.LogId, v.Type, v.Content, v.Order)
	case models.LogData:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (data, desc) VALUES (?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.Data, v.Desc)
	case models.User:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (username, password, is_admin) VALUES (?, ?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.Username, v.Password, v.IsAdmin)
	case models.Comment:
		insertQuery = fmt.Sprintf(`INSERT INTO %s (user_id, log_id, content, date) VALUES (?, ?, ?, ?)`, s.tableName)
		result, err = tx.Exec(insertQuery, v.UserId, v.LogId, v.Content, v.Date)
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
			if err := rows.Scan(&v.Id, &v.Title, &v.Date, &v.ThumbnailId, &v.Category, &v.ShortDesc); err != nil {
				return nil, err
			}
		case *models.LogItem:
			if err := rows.Scan(&v.Id, &v.LogId, &v.Type, &v.Content, &v.Order); err != nil {
				return nil, err
			}
		case *models.LogData:
			if err := rows.Scan(&v.Id, &v.Data, &v.Desc); err != nil {
				return nil, err
			}
		case *models.User:
			if err := rows.Scan(&v.Id, &v.Username, &v.Password, &v.IsAdmin); err != nil {
				return nil, err
			}
		case *models.Comment:
			if err := rows.Scan(&v.Id, &v.UserId, &v.LogId, &v.Content, &v.Date); err != nil {
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
		if err = row.Scan(&v.Id, &v.Title, &v.Date, &v.ThumbnailId, &v.Category, &v.ShortDesc); err != nil {
			return item, err
		}
	case *models.LogItem:
		if err = row.Scan(&v.Id, &v.LogId, &v.Type, &v.Content, &v.Order); err != nil {
			return item, err
		}
	case *models.LogData:
		if err := row.Scan(&v.Id, &v.Data, &v.Desc); err != nil {
			return item, err
		}
	case *models.User:
		if err := row.Scan(&v.Id, &v.Username, &v.Password, &v.IsAdmin); err != nil {
			return item, err
		}
	case *models.Comment:
		if err := row.Scan(&v.Id, &v.UserId, &v.LogId, &v.Content, &v.Date); err != nil {
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
		updateQuery := fmt.Sprintf(`UPDATE %s SET title = ?, date = ?, thumbnail_id = ?, category = ?, short_desc = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.Title, v.Date, v.ThumbnailId, v.Category, v.ShortDesc, id)
	case models.LogItem:
		updateQuery := fmt.Sprintf(`UPDATE %s SET log_id = ?, type = ?, content = ?, "order" = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.LogId, v.Type, v.Content, v.Order, id)
	case models.LogData:
		updateQuery := fmt.Sprintf(`UPDATE %s SET data = ?, desc = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.Data, v.Desc, id)
	case models.User:
		updateQuery := fmt.Sprintf(`UPDATE %s SET username = ?, password = ?, is_admin = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.Username, v.Password, v.IsAdmin, id)
	case models.Comment:
		updateQuery := fmt.Sprintf(`UPDATE %s SET user_id = ?, log_id = ?, content = ?, date = ? WHERE id = ?`, s.tableName)
		_, err = s.db.Exec(updateQuery, v.UserId, v.LogId, v.Content, v.Date, id)
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
			if err = rows.Scan(&v.Id, &v.Title, &v.Date, &v.ThumbnailId, &v.Category, &v.ShortDesc); err != nil {
				return nil, err
			}
		case *models.LogItem:
			if err = rows.Scan(&v.Id, &v.LogId, &v.Type, &v.Content, &v.Order); err != nil {
				return nil, err
			}
		case *models.LogData:
			if err := rows.Scan(&v.Id, &v.Data, &v.Desc); err != nil {
				return nil, err
			}
		case *models.User:
			if err := rows.Scan(&v.Id, &v.Username, &v.Password, &v.IsAdmin); err != nil {
				return nil, err
			}
		case *models.Comment:
			if err := rows.Scan(&v.Id, &v.UserId, &v.LogId, &v.Content, &v.Date); err != nil {
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
