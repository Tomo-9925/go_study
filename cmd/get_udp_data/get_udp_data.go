package main

import (
	"fmt"

	"github.com/tomo-9925/go_study/pkg/monitor"
)

func main() {
	data, err := monitor.GetAllUDPData()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
}
