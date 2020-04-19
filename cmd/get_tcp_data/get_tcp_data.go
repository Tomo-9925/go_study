package main

import (
	"fmt"

	"github.com/tomo-9925/go_study/pkg/monitor"
)

func main() {
	data, err := monitor.GetAllTCPData()
	if err != nil {
		fmt.Println(err)
	}
	for i, d := range data {
		fmt.Println(i, d)
	}
}
