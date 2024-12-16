package models

import "time"

type LogItemType int

const (
	Text LogItemType = iota
	Image
	Title
)

func (lt LogItemType) String() string {
	return [...]string{"Text", "Image", "Title"}[lt]
}

type LogData struct {
	Id   int    `json:"id" db:"id"`
	Data []byte `json:"data" db:"data"`
}

type LogItem struct {
	Id      int         `json:"id" db:"id"`
	LogId   int         `json:"log_id" db:"log_id"`
	Type    LogItemType `json:"type" db:"type"`
	Content string      `json:"content" db:"content"`
	Order   int         `json:"order" db:"order"`
}

type Log struct {
	Id    int       `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
	Date  time.Time `json:"date" db:"date"`
	Items []LogItem `json:"items" db:"items"`
}
