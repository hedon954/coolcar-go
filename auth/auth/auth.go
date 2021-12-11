package auth

import (
	"context"

	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
	type AuthServiceServer interface {
		Login(context.Context, *LoginRequest) (*LoginResponse, error)
	}
*/

// Auth Service, 需要实现 grpc 的 AuthServiceServer 这个接口
type Service struct {
	OpenIDResolver OpenIDResolver // OpenID 解析器
	Mongo          *dao.Mongo
	Logger         *zap.Logger // zap 包的日志工具
}

// OpenID 解析器
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// Login 登录接口
func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// 解析 OpenID
	openID, err := s.OpenIDResolver.Resolve(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openID %v: ", err)
	}

	// 根据 OpenID 获取 AccountID
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	// 日志
	s.Logger.Info("received code:", zap.String("code", req.Code))

	return &authpb.LoginResponse{
		AccessToken: "token for account id: " + accountID,
	}, nil
}
