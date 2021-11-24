package main

import (
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/wechat"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	//  zap 日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}

	// grpc 是 tcp 服务
	listenr, err := net.Listen("tcp", ":9090")
	if err != nil {
		logger.Fatal("cannot create grpc listner: %v", zap.Error(err))
	}

	// 开启一个 grpc 服务
	s := grpc.NewServer()

	// 注册 grpc 服务
	// 参数1：grpc.Server
	// 参数2：AuthServiceServer
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "wx2f9adc3f3ef8f540",
			AppSecret: "654e58d975b0fcde812beacd54f8a6c8",
		},
		Logger: logger,
	})

	err = s.Serve(listenr)
	if err != nil {
		logger.Fatal("cannot server", zap.Error(err))
	}
}

// 自定义 ZAP LOGGER
// func newZapLogger() (*zap.Logger, error) {
// 	cfg := zap.NewDevelopmentConfig()
// 	cfg.EncoderConfig.TimeKey = ""
// 	return cfg.Build()
// }
