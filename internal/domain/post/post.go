package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	UUID      string
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(title, content string) *Post {
	return &Post{
		UUID:      uuid.New().String(),
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}