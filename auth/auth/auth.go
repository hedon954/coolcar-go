package auth

import (
	"context"

	authpb "coolcar/auth/api/gen/v1"

	"go.uber.org/zap"
)

/*
	type AuthServiceServer interface {
		Login(context.Context, *LoginRequest) (*LoginResponse, error)
	}
*/

type Service struct {
	Logger *zap.Logger // zap 包的日志工具
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code:", zap.String("code", req.Code))
	return &authpb.LoginResponse{
		AccessToken: req.Code + "-accessToken",
	}, nil
}
