package postGRPC

import (
	"context"
	"errors"

	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/post"
	application "github.com/ruslanukhlin/SwiftTalk.post-service/internal/application/post"
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrShortTitle = status.Error(codes.InvalidArgument, domain.ErrShortTitle.Error())
	ErrLongTitle = status.Error(codes.InvalidArgument, domain.ErrLongTitle.Error())
	ErrShortContent = status.Error(codes.InvalidArgument, domain.ErrShortContent.Error())
	ErrLongContent = status.Error(codes.InvalidArgument, domain.ErrLongContent.Error())
	ErrPostNotFound = status.Error(codes.NotFound, domain.ErrPostNotFound.Error())
	ErrInternal = status.Error(codes.Internal, "Внутренняя ошибка сервера")
)

type PostGRPCHandler struct {
	pb.UnimplementedPostServiceServer
	postApp *application.PostApp
}

func NewPostGRPCHandler(postApp *application.PostApp) *PostGRPCHandler {
	return &PostGRPCHandler{
		postApp: postApp,
	}
}

func (h *PostGRPCHandler) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	if err := h.postApp.CreatePost(&domain.CreatePostInput{
		Title: req.Title,
		Content: req.Content,
		Images: req.Images,
	}); err != nil {
		switch err {
		case domain.ErrShortTitle:
			return nil, ErrShortTitle
		case domain.ErrLongTitle:
			return nil, ErrLongTitle
		case domain.ErrShortContent:
			return nil, ErrShortContent
		case domain.ErrLongContent:
			return nil, ErrLongContent
		default:
			return nil, ErrInternal
		}
	}

	return &pb.CreatePostResponse{}, nil
}

func (h *PostGRPCHandler) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	postsResponse, err := h.postApp.GetPosts(req.Page, req.Limit)

	if err != nil {
		return nil, ErrInternal
	}

	postsPb := make([]*pb.Post, len(postsResponse.Posts))
	for i, post := range postsResponse.Posts {
		imagesPb := getImages(post.Images)
		postsPb[i] = &pb.Post{
			Uuid: post.UUID,
			Title: post.Title.Value,
			Content: post.Content.Value,
			Images: imagesPb,
		}
	}

	return &pb.GetPostsResponse{
		Posts: postsPb,
		Total: postsResponse.Total,
		Page: postsResponse.Page,
		Limit: postsResponse.Limit,
	}, nil
}

func (h *PostGRPCHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	post, err := h.postApp.GetPostByUUID(req.Uuid)

	if err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, ErrInternal
	}

	imagesPb := getImages(post.Images)
	return &pb.GetPostResponse{
		Post: &pb.Post{
			Uuid: post.UUID,
			Title: post.Title.Value,
			Content: post.Content.Value,
			Images: imagesPb,
		},
	}, nil
}

func (h *PostGRPCHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if err := h.postApp.DeletePost(req.Uuid); err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, ErrInternal
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *PostGRPCHandler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	if err := h.postApp.UpdatePost(&domain.UpdatePostInput{
		UUID: req.Uuid,
		Title: req.Title,
		Content: req.Content,
		Images: req.Images,
		ImagesToDelete: req.DeletedImages,
	}); err != nil {
		switch err {
		case domain.ErrShortTitle:
			return nil, ErrShortTitle
		case domain.ErrLongTitle:
			return nil, ErrLongTitle
		case domain.ErrShortContent:
			return nil, ErrShortContent
		case domain.ErrLongContent:
			return nil, ErrLongContent
		default:
			return nil, ErrInternal
		}
	}

	return &pb.UpdatePostResponse{}, nil
}

func getImages(images []*domain.Image) []*pb.Image {
	imagesPb := make([]*pb.Image, len(images))
	for i, image := range images {
		imagesPb[i] = &pb.Image{
			Uuid: image.UUID,
			Url: image.URL,
		}
	}
	return imagesPb
}