package models

import "time"

type LogItemType int

const (
	Text LogItemType = iota
	Image
)

func (lt LogItemType) String() string {
	return [...]string{"Text", "Image"}[lt]
}

type LogItem struct {
	Type    LogItemType
	Content string
	Order   int
}

type Log struct {
	Date  time.Time
	Items []LogItem
}
