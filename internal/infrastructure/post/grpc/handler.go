package postGRPC

import (
	"context"

	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/post"
	application "github.com/ruslanukhlin/SwiftTalk.post-service/internal/application/post"
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrShortTitle = status.Error(codes.InvalidArgument, "Заголовок слишком короткий (минимум 3 символа)")
	ErrLongTitle = status.Error(codes.InvalidArgument, "Заголовок слишком длинный (максимум 255 символов)")
	ErrShortContent = status.Error(codes.InvalidArgument, "Содержание слишком короткое (минимум 3 символа)")
	ErrLongContent = status.Error(codes.InvalidArgument, "Содержание слишком длинное (максимум 100000 символов)")
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
	post, err := domain.NewPost(req.Title, req.Content)
	if err != nil {
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

	if err := h.postApp.CreatePost(post); err != nil {
		return nil, ErrInternal
	}

	return &pb.CreatePostResponse{}, nil
}

func (h *PostGRPCHandler) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	posts, err := h.postApp.GetPosts()

	if err != nil {
		return nil, ErrInternal
	}

	postsPb := make([]*pb.Post, len(posts))
	for i, post := range posts {
		postsPb[i] = &pb.Post{
			Uuid: post.UUID,
			Title: post.Title.Value,
			Content: post.Content.Value,
		}
	}

	return &pb.GetPostsResponse{
		Posts: postsPb,
	}, nil
}

func (h *PostGRPCHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	post, err := h.postApp.GetPostByUUID(req.Uuid)

	if err != nil {
		return nil, ErrInternal
	}

	return &pb.GetPostResponse{
		Post: &pb.Post{
			Uuid: post.UUID,
			Title: post.Title.Value,
			Content: post.Content.Value,
		},
	}, nil
}
