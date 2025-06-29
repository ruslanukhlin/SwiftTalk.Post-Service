package postGRPC

import (
	"context"

	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/post"
	application "github.com/ruslanukhlin/SwiftTalk.post-service/internal/application/post"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/auth"
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrShortTitle    = status.Error(codes.InvalidArgument, domain.ErrShortTitle.Error())
	ErrLongTitle     = status.Error(codes.InvalidArgument, domain.ErrLongTitle.Error())
	ErrShortContent  = status.Error(codes.InvalidArgument, domain.ErrShortContent.Error())
	ErrLongContent   = status.Error(codes.InvalidArgument, domain.ErrLongContent.Error())
	ErrPostNotFound  = status.Error(codes.NotFound, domain.ErrPostNotFound.Error())
	ErrUnauthorized  = status.Error(codes.Unauthenticated, auth.ErrInvalidToken.Error())
	ErrUserNotAuthor = status.Error(codes.Unauthenticated, auth.ErrUserNotAuthor.Error())
	ErrInvalidUUID   = status.Error(codes.InvalidArgument, domain.ErrInvalidUUID.Error())
	ErrInternal      = status.Error(codes.Internal, "Внутренняя ошибка сервера")
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, ErrUnauthorized
	}
	accessToken := values[0]

	if err := h.postApp.CreatePost(&domain.CreatePostInput{
		AccessToken: accessToken,
		Title:       req.Title,
		Content:     req.Content,
		Images:      req.Images,
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
		case auth.ErrInvalidToken:
			return nil, ErrUnauthorized
		case domain.ErrInvalidUUID:
			return nil, ErrInvalidUUID
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
			Uuid:     post.UUID,
			UserUuid: post.UserUUID,
			Title:    post.Title.Value,
			Content:  post.Content.Value,
			Images:   imagesPb,
		}
	}

	return &pb.GetPostsResponse{
		Posts: postsPb,
		Total: postsResponse.Total,
		Page:  postsResponse.Page,
		Limit: postsResponse.Limit,
	}, nil
}

func (h *PostGRPCHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	post, err := h.postApp.GetPostByUUID(req.Uuid)
	if err != nil {
		switch err {
		case domain.ErrPostNotFound:
			return nil, ErrPostNotFound
		case domain.ErrInvalidUUID:
			return nil, ErrInvalidUUID
		default:
			return nil, ErrInternal
		}
	}

	imagesPb := getImages(post.Images)
	return &pb.GetPostResponse{
		Post: &pb.Post{
			Uuid:     post.UUID,
			UserUuid: post.UserUUID,
			Title:    post.Title.Value,
			Content:  post.Content.Value,
			Images:   imagesPb,
		},
	}, nil
}

func (h *PostGRPCHandler) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, ErrUnauthorized
	}
	accessToken := values[0]

	if err := h.postApp.DeletePost(accessToken, req.Uuid); err != nil {
		switch err {
		case domain.ErrPostNotFound:
			return nil, ErrPostNotFound
		case auth.ErrUserNotAuthor:
			return nil, ErrUserNotAuthor
		case domain.ErrInvalidUUID:
			return nil, ErrInvalidUUID
		default:
			return nil, ErrInternal
		}
	}

	return &pb.DeletePostResponse{}, nil
}

func (h *PostGRPCHandler) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrUnauthorized
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, ErrUnauthorized
	}
	accessToken := values[0]

	if err := h.postApp.UpdatePost(&domain.UpdatePostInput{
		UUID:           req.Uuid,
		Title:          req.Title,
		Content:        req.Content,
		Images:         req.Images,
		ImagesToDelete: req.DeletedImages,
		AccessToken:    accessToken,
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
		case auth.ErrUserNotAuthor:
			return nil, ErrUserNotAuthor
		case domain.ErrInvalidUUID:
			return nil, ErrInvalidUUID
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
			Url:  image.URL,
		}
	}
	return imagesPb
}
