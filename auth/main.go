package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"

	"github.com/dgrijalva/jwt-go"
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

	// 获取 PRIVATE KEY
	pkFile, err := os.Open("auth/private.key")
	if err != nil {
		logger.Fatal("cannot open private.key", zap.Error(err))
	}
	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private.key", zap.Error(err))
	}
	pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pkBytes))
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
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
		Mongo:          dao.NewMongo(db),
		Logger:         logger,
		TokenGenerator: token.NewJWTGen("coolcar/auth", pk),
		TokenExpire:    2 * time.Hour,
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
