package clientGRPC

import (
	"context"
	"strings"

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

	response, err := c.authClient.VerifyToken(context.Background(), &pb.VerifyTokenRequest{
		AccessToken: strings.TrimPrefix(accessToken, "Bearer "),
	})
	if err != nil {
		return nil, auth.ErrVerifyToken
	}

	return &auth.VerifyTokenOutput{
		UserUUID: response.UserUuid,
	}, nil
}
