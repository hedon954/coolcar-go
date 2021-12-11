package main

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	// docker run -p 27017:27017 mongo:4

	// 获取 Docker 客户端
	c, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// 创建 Docker 容器
	// func (*client.Client).ContainerCreate(ctx context.Context,
	// config *container.Config,
	// hostConfig *container.HostConfig,
	// networkingConfig *network.NetworkingConfig,
	// platform *v1.Platform,
	// containerName string)
	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: "mongo:4",
		ExposedPorts: nat.PortSet{
			"27017/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"27017/tcp": []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// 启动容器
	err = c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// 5s 后关闭并删除容器
	err = c.ContainerStop(ctx, resp.ID, nil)
	if err != nil {
		panic(err)
	}
	err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}
}
