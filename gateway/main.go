package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot create zap logger: %v", err)
	}

	// 创建一个可以取消的 context
	c, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 创建一个 runtime mux，约束一些策略
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			protojson.MarshalOptions{
				UseEnumNumbers:  true, // enum 返回数值
				UseProtoNames:   true, // 使用 bf 原始名称
				EmitUnpopulated: true, // 输出零值
			},
			protojson.UnmarshalOptions{},
		},
	))

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "localhost:9090",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "localhost:9080",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}

	// 注册 grpc 服务到 grpc gateway
	for _, s := range serverConfig {
		err := s.registerFunc(
			c,
			mux,
			s.addr,
			[]grpc.DialOption{grpc.WithInsecure()},
		)
		if err != nil {
			logger.Sugar().Fatalf("cannot register %s service: %v", s.name, err)
		}
	}

	// 启动 grpc gateway
	logger.Sugar().Info("grpc gateway started at :9527")
	logger.Sugar().Fatal(http.ListenAndServe(":9527", mux))
}
