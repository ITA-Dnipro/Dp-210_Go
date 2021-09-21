package grpc

import (
	"context"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/server/grpc/proto"
	"github.com/ITA-Dnipro/Dp-210_Go/auth/internal/usecase"
	"net/http"
)

type Auth interface {
	ValidateToken(t usecase.JwtToken) (usecase.UserAuth, error)
}

type Server struct {
	proto.UnimplementedTokenValidatorServer
	auth Auth
}

func NewGrpcServer(auth Auth) *Server {
	return &Server{auth: auth}
}

func (s *Server) Validate(ctx context.Context, in *proto.Token) (*proto.ValidatedToken, error) {
	t := in.Token

	user, err := s.auth.ValidateToken(usecase.JwtToken(t))
	if err != nil {
		return &proto.ValidatedToken{StatusCode: http.StatusUnauthorized, UserId: "", UserRole: ""}, nil
	}

	return &proto.ValidatedToken{StatusCode: http.StatusOK, UserId: user.Id, UserRole: string(user.Role)}, nil
}
