package models

import (
	storeM "LogC/internal/models/store"
	"time"
)

// Request format

type CommentRequest struct {
	UserId  int    `json:"user_id"`
	LogId   int    `json:"log_id"`
	Content string `json:"content"`
}

func FromCommentRequest(comment CommentRequest) storeM.Comment {
	return storeM.Comment{
		UserId:  comment.UserId,
		LogId:   comment.LogId,
		Content: comment.Content,
		Date:    time.Now(),
	}
}

// Response format

type CommentResponse struct {
	Id       int       `json:"id"`
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
}

func ToCommentResponse(comment storeM.Comment, username string) CommentResponse {
	return CommentResponse{
		Id:       comment.Id,
		Username: username,
		Content:  comment.Content,
		Date:     comment.Date,
	}
}
