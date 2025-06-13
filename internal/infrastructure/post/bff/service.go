package bff

import (
	"context"
	"io"
	"mime/multipart"

	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/post"
)

type PostService struct {
	client pb.PostServiceClient
}

func NewPostService(client pb.PostServiceClient) *PostService {
	return &PostService{
		client: client,
	}
}

func (s *PostService) GetPosts(ctx context.Context) ([]*Post, error) {
	response, err := s.client.GetPosts(ctx, &pb.GetPostsRequest{})
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, len(response.Posts))
	for i, post := range response.Posts {
		images := post.Images
		if images == nil {
			images = make([]string, 0)
		}
		posts[i] = &Post{
			Uuid:    post.Uuid,
			Title:   post.Title,
			Content: post.Content,
			Images:  images,
		}
	}

	return posts, nil
}

func (s *PostService) CreatePost(ctx context.Context, title, content string, images []*multipart.FileHeader) error{
	imageBytes := make([][]byte, len(images))
	for i, image := range images {
		image, err := image.Open()
		if err != nil {
			return err
		}

		imageBytes[i], err = io.ReadAll(image)
		if err != nil {
			return err
		}

		defer image.Close()
	}

	_, err := s.client.CreatePost(ctx, &pb.CreatePostRequest{
		Title:   title,
		Content: content,
		Images:  imageBytes,
	})
	if err != nil {
		return err
	}

	return nil
}