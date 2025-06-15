package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	UUID      string
	UserUUID  string
	Title     Title
	Content   Content
	Images    []*Image
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPost(userUUID, title, content string, images []*Image) (*Post, error) {
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
		UserUUID:  userUUID,
		Title:     *titleValid,
		Content:   *contentValid,
		Images:    images,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
