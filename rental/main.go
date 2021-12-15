package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	//  zap 日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("cannot create zap logger", zap.Error(err))
	}

	// grpc 是 tcp 服务，需要监听端口
	listener, err := net.Listen("tcp", ":9080")
	if err != nil {
		logger.Fatal("cannot create grpc listner", zap.Error(err))
	}

	// 创建一个新的 grpc 服务
	s := grpc.NewServer()
	rentalpb.RegisterTripServiceServer(s, &trip.Service{
		Logger: logger,
	})

	// 启动服务
	err = s.Serve(listener)
	if err != nil {
		logger.Fatal("cannot server", zap.Error(err))
	}

}
