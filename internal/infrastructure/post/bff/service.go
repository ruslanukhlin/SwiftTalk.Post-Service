package bff

import (
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
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

func (s *PostService) GetPost(c *fiber.Ctx, postId string) (*Post, error) {
	payload, err := s.client.GetPost(c.Context(), &pb.GetPostRequest{
		Uuid: postId,
	})
	if err != nil {
		return nil, err
	}

	images := getImages(payload.Post.Images)
	return &Post{
		Uuid: payload.Post.Uuid,
		Title: payload.Post.Title,
		Content: payload.Post.Content,
		Images: images,
	}, nil
}

func (s *PostService) GetPosts(c *fiber.Ctx, page, limit int64) (*GetPostsResponse, error) {
	response, err := s.client.GetPosts(c.Context(), &pb.GetPostsRequest{
		Page: page,
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	posts := make([]*Post, len(response.Posts))
	for i, post := range response.Posts {
		images := getImages(post.Images)
		posts[i] = &Post{
			Uuid:    post.Uuid,
			Title:   post.Title,
			Content: post.Content,
			Images:  images,
		}
	}

	return &GetPostsResponse{
		Posts: posts,
		Total: response.Total,
		Page: response.Page,
		Limit: response.Limit,
	}, nil
}

func (s *PostService) CreatePost(c *fiber.Ctx, title, content string, images []*multipart.FileHeader) error{
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

	_, err := s.client.CreatePost(c.Context(), &pb.CreatePostRequest{
		Title:   title,
		Content: content,
		Images:  imageBytes,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) UpdatePost(c *fiber.Ctx, postId, title, content string, images []*multipart.FileHeader, deletedImages []string) error {
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

	_, err := s.client.UpdatePost(c.Context(), &pb.UpdatePostRequest{
		Uuid: postId,
		Title: title,
		Content: content,
		Images: imageBytes,
		DeletedImages: deletedImages,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *PostService) DeletePost(c *fiber.Ctx, postId string) error {
	_, err := s.client.DeletePost(c.Context(), &pb.DeletePostRequest{
		Uuid: postId,
	})
	if err != nil {
		return err
	}

	return nil
}

func getImages(images []*pb.Image) []*Image {
	imagesBff := make([]*Image, len(images))
	for i, image := range images {
		imagesBff[i] = &Image{
			Uuid: image.Uuid,
			Url: image.Url,
		}
	}
	return imagesBff
}