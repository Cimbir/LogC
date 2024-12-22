package models

import "time"

type Comment struct {
	Id      int       `json:"id" db:"id"`
	UserId  int       `json:"user_id" db:"user_id"`
	LogId   int       `json:"log_id" db:"log_id"`
	Content string    `json:"content" db:"content"`
	Date    time.Time `json:"date" db:"date"`
}
