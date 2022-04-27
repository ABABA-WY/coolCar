package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const containerPort = "27017/tcp"

func main() {
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
	//关闭
	defer func() {
		err = newEnvClient.ContainerRemove(c, create.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			panic(err)
		}
	}()
	//开启
	err = newEnvClient.ContainerStart(c, create.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	//查看监听的port
	inspectRes, err := newEnvClient.ContainerInspect(c, create.ID)
	if err != nil {
		panic(err)
	}
	hostPort := inspectRes.NetworkSettings.Ports[containerPort][0]
	//告知外界MongoDB启动的URI
	fmt.Printf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)

}
