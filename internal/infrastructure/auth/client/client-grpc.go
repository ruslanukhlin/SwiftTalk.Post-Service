package clientGRPC

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	pb "github.com/ruslanukhlin/SwiftTalk.Common/gen/auth"
	"github.com/ruslanukhlin/SwiftTalk.post-service/internal/domain/auth"
)

var _ auth.AuthRepository = &ClientGRPC{}

type ClientGRPC struct {
	authClient pb.AuthServiceClient
}

func NewClientGRPC(authClient pb.AuthServiceClient) *ClientGRPC {
	return &ClientGRPC{authClient: authClient}
}

func (c *ClientGRPC) VerifyToken(accessToken string) (*auth.VerifyTokenOutput, error) {
	if accessToken == "" {
		return nil, auth.ErrInvalidToken
	}

	md := metadata.New(map[string]string{"authorization": accessToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	response, err := c.authClient.VerifyToken(ctx, &pb.VerifyTokenRequest{})
	if err != nil {
		log.Println("VerifyToken error:", err)
		return nil, auth.ErrVerifyToken
	}

	return &auth.VerifyTokenOutput{
		UserUUID: response.UserUuid,
	}, nil
}
