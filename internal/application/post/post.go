package application

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/google/uuid"
	s3 "github.com/ruslanukhlin/SwiftTalk.Common/core/s3"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/auth"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"github.com/ruslanukhlin/SwiftTalk.post-service/pkg/config"
)

var _ post.PostService = &PostApp{}

type PostApp struct {
	PostRepository post.PostRepository
	authClient     auth.AuthRepository
	s3             *s3.S3
	cfg            *config.Config
}

func NewPostApp(postRepo post.PostRepository, s3 *s3.S3, authClient auth.AuthRepository) *PostApp {
	cfg := config.LoadConfigFromEnv()
	return &PostApp{
		PostRepository: postRepo,
		authClient:     authClient,
		s3:             s3,
		cfg:            cfg,
	}
}

func (a *PostApp) CreatePost(input *post.CreatePostInput) error {
	verifyTokenOutput, err := a.authClient.VerifyToken(input.AccessToken)
	if err != nil {
		return err
	}

	images := a.getImages(input.Images)

	err = a.s3.UploadFiles(context.Background(), images.Readers, images.Urls)
	if err != nil {
		log.Printf("Ошибка при загрузке изображений: %v", err)
		return err
	}

	post, err := post.NewPost(verifyTokenOutput.UserUUID, input.Title, input.Content, images.Domain)
	if err != nil {
		return err
	}

	return a.PostRepository.Save(post)
}

func (a *PostApp) GetPosts(page, limit int64) (*post.GetPostsResponse, error) {
	return a.PostRepository.FindAll(page, limit)
}

func (a *PostApp) GetPostByUUID(uuid string) (*post.Post, error) {
	return a.PostRepository.FindByUUID(uuid)
}

func (a *PostApp) UpdatePost(input *post.UpdatePostInput) error {
	verifyTokenOutput, err := a.authClient.VerifyToken(input.AccessToken)
	if err != nil {
		return err
	}

	foundPost, err := a.PostRepository.FindByUUID(input.UUID)
	if err != nil {
		return err
	}

	if verifyTokenOutput.UserUUID != foundPost.UserUUID {
		return auth.ErrUserNotAuthor
	}

	images := a.getImages(input.Images)

	err = a.s3.UploadFiles(context.Background(), images.Readers, images.Urls)
	if err != nil {
		return err
	}

	err = a.PostRepository.DeleteImages(foundPost.UUID, input.ImagesToDelete)
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

	foundPost.Title = post.Title{Value: input.Title}
	foundPost.Content = post.Content{Value: input.Content}
	foundPost.Images = images.Domain

	return a.PostRepository.Update(foundPost)
}

func (a *PostApp) DeletePost(accessToken, uuid string) error {
	verifyTokenOutput, err := a.authClient.VerifyToken(accessToken)
	if err != nil {
		return err
	}

	post, err := a.PostRepository.FindByUUID(uuid)
	if err != nil {
		return err
	}

	if verifyTokenOutput.UserUUID != post.UserUUID {
		return auth.ErrUserNotAuthor
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
	Urls    []string
	Domain  []*post.Image
}

func (a *PostApp) getImages(images [][]byte) *Images {
	imagesReaders := make([]io.Reader, len(images))
	imagesUrls := make([]string, len(images))
	imagesDomain := make([]*post.Image, len(images))
	for i, image := range images {
		newUuid := uuid.New().String()
		pathToS3 := a.cfg.S3.BucketUrl + newUuid
		pathToBucket := a.cfg.S3.BucketFolder + "/" + newUuid

		imagesUrls[i] = pathToBucket
		imagesDomain[i] = post.NewImage(newUuid, pathToS3)
		imagesReaders[i] = bytes.NewReader(image)
	}
	return &Images{
		Readers: imagesReaders,
		Urls:    imagesUrls,
		Domain:  imagesDomain,
	}
}
