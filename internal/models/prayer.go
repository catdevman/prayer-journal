package models

import (
	"time"
)

type Prayer struct {
	ID         string    `json:"id" dynamodbav:"id"`
	Title      string    `json:"title" dynamodbav:"title"`
	Content    string    `json:"content" dynamodbav:"content"`
	IsAnswered bool      `json:"is_answered" dynamodbav:"is_answered"`
	CreatedAt  time.Time `json:"created_at" dynamodbav:"created_at"`
}
