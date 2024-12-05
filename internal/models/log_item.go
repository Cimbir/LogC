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
	Id   int
	Data []byte
}

type LogItem struct {
	Id      int
	LogId   int
	Type    LogItemType
	Content string
	Order   int
}

type Log struct {
	Id    int
	Title string
	Date  time.Time
	Items []LogItem
}
