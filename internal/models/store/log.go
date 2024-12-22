package models

import "time"

// LogItemType is an enum for the type of log item
type LogItemType int

const (
	Text LogItemType = iota
	Image
	Title
	Quote
)

func (lt LogItemType) String() string {
	return [...]string{"Text", "Image", "Title", "Quote"}[lt]
}

func GetLogItemType(lts string) LogItemType {
	switch lts {
	case "Text":
		return Text
	case "Image":
		return Image
	case "Title":
		return Title
	case "Quote":
		return Quote
	default:
		return -1
	}
}

// LogCategory is an enum for the category of log
type LogCategory int

const (
	Other LogCategory = iota
	Tech
	Art
	Review
)

func (lc LogCategory) String() string {
	return [...]string{"Other", "Tech", "Art", "Review"}[lc]
}

func GetLogCategory(lcs string) LogCategory {
	switch lcs {
	case "Other":
		return Other
	case "Tech":
		return Tech
	case "Art":
		return Art
	case "Review":
		return Review
	default:
		return -1
	}
}

type LogData struct {
	Id   int    `json:"id" db:"id"`
	Data []byte `json:"data" db:"data"`
	Desc string `json:"desc" db:"desc"`
}

type LogItem struct {
	Id      int         `json:"id" db:"id"`
	LogId   int         `json:"log_id" db:"log_id"`
	Type    LogItemType `json:"type" db:"type"`
	Content string      `json:"content" db:"content"`
	Order   int         `json:"order" db:"order"`
}

type Log struct {
	Id          int         `json:"id" db:"id"`
	Title       string      `json:"title" db:"title"`
	Date        time.Time   `json:"date" db:"date"`
	ThumbnailId int         `json:"thumbnail_id" db:"thumbnail_id"`
	Category    LogCategory `json:"category" db:"category"`
	ShortDesc   string      `json:"short_desc" db:"short_desc"`
}
