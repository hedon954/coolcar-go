package shared_server

import (
	shared_auth "coolcar/shared/auth"

	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Name                  string
	Address               string
	AuthPublicKeyFilePath string
	RegisterFunc          func(*grpc.Server)
	Logger                *zap.Logger
}

// RunGRPCServer runs a grpc server
func RunGRPCServer(config *GRPCConfig) error {
	nameField := zap.String("name", config.Name)

	// grpc 是 tcp 服务，需要监听端口
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		config.Logger.Fatal("cannot create grpc listner", nameField, zap.Error(err))
	}

	var opts []grpc.ServerOption

	// auth service 不需要这个拦截器
	if config.AuthPublicKeyFilePath != "" {
		// 获取登录拦截器
		in, err := shared_auth.Interceptor(config.AuthPublicKeyFilePath)
		if err != nil {
			config.Logger.Fatal("cannot get auth interceptor", nameField, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	// 创建一个新的 grpc 服务
	s := grpc.NewServer(opts...)
	config.RegisterFunc(s)

	// 启动服务
	config.Logger.Sugar().Infof("%s service started at %s", config.Name, config.Address)
	return s.Serve(listener)
}
