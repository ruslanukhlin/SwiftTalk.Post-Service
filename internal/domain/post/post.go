package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	UUID      string
	Title     Title
	Content   Content
	Images    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(title, content string, images []string) (*Post, error) {
	titleValid, err := NewTitle(title)
	if err != nil {
		return nil, err
	}

	contentValid, err := NewContent(content)
	if err != nil {
		return nil, err
	}

	return &Post{
		UUID:      uuid.New().String(),
		Title:     *titleValid,
		Content:   *contentValid,
		Images:    images,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}