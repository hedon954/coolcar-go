package auth

import (
	"context"
	"time"

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
	Logger         *zap.Logger    // zap 包的日志工具
	TokenGenerator TokenGenerator // token 生成器
	TokenExpire    time.Duration  // token 过期时间
}

// OpenID 解析器
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// TokenGenerator 为 accountID 生成一个 token
type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
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

	// 生成 token
	token, err := s.TokenGenerator.GenerateToken(accountID.String(), s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
