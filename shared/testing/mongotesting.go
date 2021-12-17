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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_IMAGE  = "mongo:4"
	CONTAINER_PORT = "27017/tcp"
)

var mongoURI string

const defaultMongoURI = "mongodb://localhost:27017"

// RunWithMongoInDocker runs the tests with
// a mongodb instance in a docker container
func RunWithMongoInDocker(m *testing.M) int {
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
	mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

	return m.Run()
}

// NewClient
func NewClient(c context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("mongo uri no set, please run RunWithMongoInDocker in TestMain")
	}
	return mongo.Connect(c, options.Client().ApplyURI(mongoURI))
}

// NewDefaultClient
func NewDefaultClient(c context.Context) (*mongo.Client, error) {
	return mongo.Connect(c, options.Client().ApplyURI(defaultMongoURI))
}

// SetupIndexes setups indexes for mongodb
func SetupIndexes(c context.Context, d *mongo.Database) error {
	_, err := d.Collection("account").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{
				Key:   "open_id",
				Value: 1,
			},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = d.Collection("trip").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "trip.accountid", Value: 1},
			{Key: "trip.status", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
			"trip.status": 1,
		}),
	})

	return err
}
