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

type LogItem struct {
	Type    LogItemType
	Content string
	Order   int
}

type Log struct {
	Title string
	Date  time.Time
	Items []LogItem
}
