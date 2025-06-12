package bff

import (
	"context"

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
		posts[i] = &Post{
			Uuid:    post.Uuid,
			Title:   post.Title,
			Content: post.Content,
		}
	}

	return posts, nil
}

func (s *PostService) CreatePost(ctx context.Context, payload *CreatePostPayload) error{
	_, err := s.client.CreatePost(ctx, &pb.CreatePostRequest{
		Title:   payload.Title,
		Content: payload.Content,
	})
	if err != nil {
		return err
	}

	return nil
}