package mongotesting

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	MONGODB_IMAGE  = "mongo:4"
	CONTAINER_PORT = "27017/tcp"
)

// RunWithMongoInDocker runs the tests with
// a mongodb instance in a docker container
func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
	// docker run -p 27017:27017 mongo:4

	// 获取 Docker 客户端
	c, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	// 创建 Docker 容器
	// func (*client.Client).ContainerCreate(ctx context.Context,
	// 		config *container.Config,
	// 		hostConfig *container.HostConfig,
	// 		networkingConfig *network.NetworkingConfig,
	// 		platform *v1.Platform,
	// 		containerName string)
	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: MONGODB_IMAGE,
		ExposedPorts: nat.PortSet{
			CONTAINER_PORT: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			CONTAINER_PORT: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				},
			},
		},
	}, nil, nil, "coolcar-mongo")
	if err != nil {
		panic(err)
	}

	// 获取容器 ID
	containerID := resp.ID

	// 启动容器
	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}

	// 关闭容器
	defer func() {
		time.Sleep(5 * time.Second)
		err = c.ContainerStop(ctx, containerID, nil)
		if err != nil {
			panic(err)
		}
		err = c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()

	// 获取主机地址和端口号
	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		log.Fatal(err)
	}
	hostPort := inspRes.NetworkSettings.Ports["27017/tcp"][0]
	*mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}
