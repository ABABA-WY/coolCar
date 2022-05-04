package testing

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

const containerPort = "27017/tcp"

var mongoURI string

const defaultMongoURI = "mongodb://admin:123456@localhost:27017"

//RunWithMongoInDocker run the test with a mongodb instance in a docker container
func RunWithMongoInDocker(m *testing.M) int {
	newEnvClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	c := context.Background()
	create, err := newEnvClient.ContainerCreate(c, &container.Config{
		Image: "mongo",
		ExposedPorts: nat.PortSet{
			containerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "127.0.0.1",
					HostPort: "0", //自动挑
				},
			},
		},
	}, nil, nil, "",
	)
	if err != nil {
		panic(err)
	}
	containerID := create.ID
	//关闭
	defer func() {
		err = newEnvClient.ContainerRemove(c, containerID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()
	//开启
	err = newEnvClient.ContainerStart(c, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	//查看监听的port
	inspectRes, err := newEnvClient.ContainerInspect(c, containerID)
	if err != nil {
		panic(err)
	}
	hostPort := inspectRes.NetworkSettings.Ports[containerPort][0]
	mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)
	fmt.Println(mongoURI)
	return m.Run()
}
func NewMongoClient(c context.Context) (*mongo.Client, error) {
	if mongoURI == "" {
		return nil, fmt.Errorf("invalid mongoURI\n")
	}
	return mongo.Connect(c, options.Client().ApplyURI(mongoURI))
}
func NewDefaultMongoClient(c context.Context) (*mongo.Client, error) {
	return mongo.Connect(c, options.Client().ApplyURI(defaultMongoURI))
}

// SetupIndexes 创建索引
func SetupIndexes(c context.Context, d *mongo.Database) error {
	//在account表为openid创建唯一索引
	_, err := d.Collection("account").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "open_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	_, err = d.Collection("trip").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{ //有序的
			{Key: "trip.accountid", Value: 1},
			{Key: "trip.tripstatus", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{
			"trip.tripstatus": 1,
		}),
	})

	_, err = d.Collection("profile").Indexes().CreateOne(c, mongo.IndexModel{
		Keys: bson.D{
			{Key: "accountid", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	return err
}
