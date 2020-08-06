package monitor

import (
	"context"
	"net"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var containerInformations []*Container

// Container はDocker Engine APIから取得できるデータの構造体です．
type Container struct {
	ID        string
	IP        net.IP
	Name      string
	NetworkID string
	Pid       int
}

// GetContainerInformations はDocker Engine APIを使用して必要なコンテナの情報を取得します．
func GetContainerInformations() ([]*Container, error) {
	var containerInformations []*Container

	cli, err := client.NewClientWithOpts(client.WithVersion(APIVersion)) // APIに接続
	if err != nil {
		return containerInformations, err
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{}) // コンテナのリストを取得
	if err != nil {
		return containerInformations, err
	}

	for _, container := range containers {
		var c Container
		inspect, err := cli.ContainerInspect(context.Background(), container.ID)
		if err != nil {
			return containerInformations, err
		}

		c.ID = inspect.ID
		c.IP = net.ParseIP(inspect.NetworkSettings.IPAddress)
		c.Name = inspect.Name
		c.Pid = inspect.State.Pid

		containerInformations = append(containerInformations, &c)
	}

	return containerInformations, nil
}

func storeContainerInfo() error {
	containerInfo, err := GetContainerInformations()
	containerInformations = containerInfo
	if err != nil {
		return nil
	}
	return nil
}
