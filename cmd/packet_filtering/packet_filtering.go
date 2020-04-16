package main

import (
	"fmt"

	"github.com/tomo-9925/go_study/pkg/monitor"
)

// 定数
// const (
// 	Debug             = true // デバッグ機能
// 	QueueNum          = 2    // キュー番号
// 	MaxPacketsInQueue = 100  // キューのサイズ
// )

func main() {
	data, err := monitor.GetAllTCPData()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*data)
}
