package postGRPC

import (
	"context"

	pb "github.com/ruslanukhlin/SwiftTalk.common/gen/post"
	application "github.com/ruslanukhlin/SwiftTalk.post-service/internal/application/post"
	domain "github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/post"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	post := domain.NewPost(req.Title, req.Content);

	if err := h.postApp.CreatePost(post); err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при создании поста: %v", err)
	}

	return &pb.CreatePostResponse{}, nil
}

func (h *PostGRPCHandler) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	posts, err := h.postApp.GetPosts()

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении постов: %v", err)
	}

	postsPb := make([]*pb.Post, len(posts))
	for i, post := range posts {
		postsPb[i] = &pb.Post{
			Uuid: post.UUID,
			Title: post.Title,
			Content: post.Content,
		}
	}

	return &pb.GetPostsResponse{
		Posts: postsPb,
	}, nil
}

func (h *PostGRPCHandler) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	post, err := h.postApp.GetPostByUUID(req.Uuid)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка при получении поста: %v", err)
	}

	return &pb.GetPostResponse{
		Post: &pb.Post{
			Uuid: post.UUID,
			Title: post.Title,
			Content: post.Content,
		},
	}, nil
}
