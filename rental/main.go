package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	shared_server "coolcar/shared/server"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	//  zap 日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("cannot create zap logger", zap.Error(err))
	}

	// 启动一个 grpc 服务
	err = shared_server.RunGRPCServer(&shared_server.GRPCConfig{
		Name:                  "rental ",
		Address:               ":9080",
		AuthPublicKeyFilePath: "shared/auth/public.key",
		Logger:                logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
	})

	logger.Sugar().Fatal(err)
}
