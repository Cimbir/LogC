package models

import (
	storeM "LogC/internal/models/store"
	"time"
)

// Request formats

type LogItemRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func FromLogItemRequest(item LogItemRequest, order int, logId int) storeM.LogItem {
	return storeM.LogItem{
		Type:    storeM.GetLogItemType(item.Type),
		Content: item.Content,
		Order:   order,
		LogId:   logId,
	}
}

type LogRequest struct {
	Title       string           `json:"title"`
	ThumbnailId int              `json:"thumbnail_id"`
	Category    string           `json:"category"`
	ShortDesc   string           `json:"short_desc"`
	Items       []LogItemRequest `json:"items"`
}

func FromLogRequest(log LogRequest) storeM.Log {
	return storeM.Log{
		Title:       log.Title,
		Date:        time.Now(),
		ThumbnailId: log.ThumbnailId,
		Category:    storeM.GetLogCategory(log.Category),
		ShortDesc:   log.ShortDesc,
	}
}

// Response formats

type LogItemResponse struct {
	Id      int    `json:"id"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Order   int    `json:"order"`
}

func ToLogItemResponse(item storeM.LogItem) LogItemResponse {
	return LogItemResponse{
		Id:      item.Id,
		Type:    item.Type.String(),
		Content: item.Content,
		Order:   item.Order,
	}
}

type LogResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	ThumbnailId int       `json:"thumbnail_id"`
	Category    string    `json:"category"`
	ShortDesc   string    `json:"short_desc"`
}

func ToLogResponse(log storeM.Log) LogResponse {
	// Create a new response
	response := LogResponse{
		Id:          log.Id,
		Title:       log.Title,
		Date:        log.Date,
		ThumbnailId: log.ThumbnailId,
		Category:    log.Category.String(),
		ShortDesc:   log.ShortDesc,
	}
	return response
}

type FullLogResponse struct {
	Id          int               `json:"id"`
	Title       string            `json:"title"`
	Date        time.Time         `json:"date"`
	ThumbnailId int               `json:"thumbnail_id"`
	Category    string            `json:"category"`
	ShortDesc   string            `json:"short_desc"`
	Items       []LogItemResponse `json:"items"`
	Comments    []CommentResponse `json:"comments"`
}

func ToFullLogResponse(log storeM.Log, items []storeM.LogItem, comments []storeM.Comment, usernames map[int]string) FullLogResponse {
	itemsResponse := []LogItemResponse{}
	for _, item := range items {
		itemsResponse = append(itemsResponse, ToLogItemResponse(item))
	}

	commentsResponse := []CommentResponse{}
	for _, comment := range comments {
		commentsResponse = append(commentsResponse, ToCommentResponse(comment, usernames[comment.UserId]))
	}

	// Create a new response
	response := FullLogResponse{
		Id:          log.Id,
		Title:       log.Title,
		Date:        log.Date,
		ThumbnailId: log.ThumbnailId,
		Category:    log.Category.String(),
		ShortDesc:   log.ShortDesc,
		Items:       itemsResponse,
		Comments:    commentsResponse,
	}
	return response
}

type LogByCategoryResponse struct {
	Category map[string][]LogResponse `json:"category"`
}

func ToLogByCategoryResponse(logs []storeM.Log, amount int) LogByCategoryResponse {
	response := LogByCategoryResponse{
		Category: map[string][]LogResponse{},
	}
	for _, log := range logs {
		if _, ok := response.Category[log.Category.String()]; !ok {
			response.Category[log.Category.String()] = []LogResponse{}
		}
		if len(response.Category[log.Category.String()]) >= amount {
			continue
		}
		response.Category[log.Category.String()] = append(response.Category[log.Category.String()], ToLogResponse(log))
	}
	return response
}
