package application

import (
	"bytes"
	"context"
	"io"

	"github.com/google/uuid"
	s3 "github.com/ruslanukhlin/SwiftTalk.common/core/s3"
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
)

var _ domain.PostService = &PostApp{}

type PostApp struct {
	PostRepository domain.PostRepository
	s3             *s3.S3
	cfg            *config.Config
}

func NewPostApp(postRepo domain.PostRepository, s3 *s3.S3) *PostApp {
	cfg := config.LoadConfigFromEnv()
	return &PostApp{
		PostRepository: postRepo,
		s3:             s3,
		cfg:            cfg,
	}
}

func (a *PostApp) CreatePost(input *domain.CreatePostInput) error {
	images := a.getImages(input.Images)

	err := a.s3.UploadFiles(context.Background(), images.Readers, images.Urls)
	if err != nil {
		return err
	}

	post, err := domain.NewPost(input.Title, input.Content, images.Domain)
	if err != nil {
		return err
	}

	return a.PostRepository.Save(post)
}

func (a *PostApp) GetPosts(page, limit int64) (*domain.GetPostsResponse, error) {
	return a.PostRepository.FindAll(page, limit)
}

func (a *PostApp) GetPostByUUID(uuid string) (*domain.Post, error) {
	return a.PostRepository.FindByUUID(uuid)
}

func (a *PostApp) UpdatePost(input *domain.UpdatePostInput) error {
	post, err := a.PostRepository.FindByUUID(input.UUID)
	if err != nil {
		return err
	}

	images := a.getImages(input.Images)

	err = a.s3.UploadFiles(context.Background(), images.Readers, images.Urls)
	if err != nil {
		return err
	}

	err = a.PostRepository.DeleteImages(post.UUID, input.ImagesToDelete)
	if err != nil {
		return err
	}

	imageS3Deletes := make([]string, len(input.ImagesToDelete))
	for i, image := range input.ImagesToDelete {
		imageS3Deletes[i] = a.cfg.S3.BucketFolder + "/" + image
	}

	err = a.s3.DeleteFiles(context.Background(), imageS3Deletes)
	if err != nil {
		return err
	}

	post.Title = domain.Title{Value: input.Title}
	post.Content = domain.Content{Value: input.Content}
	post.Images = images.Domain

	return a.PostRepository.Update(post)
}

func (a *PostApp) DeletePost(uuid string) error {
	post, err := a.PostRepository.FindByUUID(uuid)
	if err != nil {
		return err
	}

	s3Keys := make([]string, len(post.Images))
	for i, imageUrl := range post.Images {
		s3Keys[i] = a.cfg.S3.BucketFolder + "/" + imageUrl.UUID
	}

	err = a.s3.DeleteFiles(context.Background(), s3Keys)
	if err != nil {
		return err
	}
	return a.PostRepository.Delete(uuid)
}

type Images struct {
	Readers []io.Reader
	Uuids   []string
	Urls    []string
	Domain  []*domain.Image
}

func (a *PostApp) getImages(images [][]byte) *Images {
	imagesReaders := make([]io.Reader, len(images))
	imagesUuids := make([]string, len(images))
	imagesUrls := make([]string, len(images))
	imagesDomain := make([]*domain.Image, len(images))
	for i, image := range images {
		imagesUuids[i] = uuid.New().String()
		imagesUrls[i] = a.cfg.S3.BucketUrl + "/" + imagesUuids[i]
		imagesDomain[i] = domain.NewImage(imagesUuids[i], a.cfg.S3.BucketUrl+"/"+imagesUrls[i])
		imagesReaders[i] = bytes.NewReader(image)
	}
	return &Images{
		Readers: imagesReaders,
		Uuids:   imagesUuids,
		Urls:    imagesUrls,
		Domain:  imagesDomain,
	}
}
