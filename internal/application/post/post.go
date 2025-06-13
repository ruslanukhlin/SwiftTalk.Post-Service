package application

import (
	"bytes"
	"context"

	"github.com/google/uuid"
	s3 "github.com/ruslanukhlin/SwiftTalk.common/core/s3"
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
)

var _ domain.PostService = &PostApp{}

type PostApp struct {
	domain.PostRepository
	s3 *s3.S3
	cfg *config.Config
}

func NewPostApp(postRepo domain.PostRepository, s3 *s3.S3, cfg *config.Config) *PostApp {
	return &PostApp{
		PostRepository: postRepo,
		s3: s3,
		cfg: cfg,
	}
}

func (a *PostApp) CreatePost(title, content string, images [][]byte) error {
	uuids := make([]string, len(images))
	for i, image := range images {
		uuids[i] = uuid.New().String()
		err := a.s3.UploadFile(context.Background(), bytes.NewReader(image), "posts/" + uuids[i])
		if err != nil {
			return err
		}
	}

	imagesUrl := make([]string, len(uuids))
	for i, uuid := range uuids {
		imagesUrl[i] = a.cfg.S3.BucketUrl + "/posts/" + uuid
	}

	post, err := domain.NewPost(title, content, imagesUrl)
	if err != nil {
		return err
	}

	return a.PostRepository.Save(post)
}

func (a *PostApp) GetPosts() ([]*domain.Post, error) {
	return a.PostRepository.FindAll()
}

func (a *PostApp) GetPostByUUID(uuid string) (*domain.Post, error) {
	return a.PostRepository.FindByUUID(uuid)
}

func (a *PostApp) UpdatePost(post *domain.Post) error {
	return a.PostRepository.Update(post)
}

func (a *PostApp) DeletePost(uuid string) error {
	return a.PostRepository.Delete(uuid)
}