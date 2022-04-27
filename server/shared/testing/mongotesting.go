package mongotesting

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"testing"
)

const containerPort = "27017/tcp"

//RunWithMongoInDocker run the test with a mongodb instance in a docker container
func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
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
	*mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)
	return m.Run()
}
