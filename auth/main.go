package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/wechat"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		logger.Fatal("cannot create grpc listner", zap.Error(err))
	}

	// 获取一个 MongoDB Client 对象
	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017").SetRetryWrites(false))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	db := mongoClient.Database("coolcar")

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
		Mongo:  dao.NewMongo(db),
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
