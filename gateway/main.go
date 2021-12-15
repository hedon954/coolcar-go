package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
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

	// 注册 auth 服务
	err := authpb.RegisterAuthServiceHandlerFromEndpoint(
		c,
		mux,
		"localhost:9090",
		[]grpc.DialOption{
			grpc.WithInsecure(),
		},
	)
	if err != nil {
		log.Fatalf("cannot register auth service: %v", err)
	}

	// 注册 rental 服务
	err = rentalpb.RegisterTripServiceHandlerFromEndpoint(
		c,
		mux,
		"localhost:9080",
		[]grpc.DialOption{
			grpc.WithInsecure(),
		},
	)
	if err != nil {
		log.Fatalf("cannot register rental service: %v", err)
	}

	// 启动 grpc gateway
	log.Fatal(http.ListenAndServe(":9527", mux))
}
