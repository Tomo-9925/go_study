package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	APIVersion = "1.40"
)

func main() {
	// APIのバージョンは環境に合わせて変更する
	cli, err := client.NewClientWithOpts(client.WithVersion(APIVersion))
	if err != nil {
		panic(err)
	}

	// 実行中のコンテナの情報を取得
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		// https://godoc.org/github.com/docker/docker/client
		fmt.Printf("Container:\n%v\n", container)
		fmt.Printf("State: %v\n", container.State)
		fmt.Printf("Status: %v\n\n", container.Status)
		inspect, _ := cli.ContainerInspect(context.Background(), container.ID)
		fmt.Printf("Inspect:\n%v\n", inspect)
		fmt.Printf("State: %v\n", inspect.State)
		fmt.Printf("Pid: %v\n", inspect.State.Pid)
		fmt.Printf("NetworkMode: %v\n", inspect.HostConfig.NetworkMode)
		fmt.Printf("PortBindings: %v\n\n", inspect.HostConfig.PortBindings)
		fmt.Printf("NetworkSettings:\n%v\n", inspect.NetworkSettings)
		fmt.Printf("IPAddress: %v\n\n", inspect.NetworkSettings.IPAddress)
	}
}
