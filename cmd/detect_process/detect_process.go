package main

import (
	"fmt"
	"os"

	"github.com/tomo-9925/go_study/pkg/monitor"
)

func main() {
	tcp, err := monitor.GetAllTCPData()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(*tcp)
	var num uint16
	fmt.Print("どのTCPの情報を取得しますか．配列番号を指定してください．\n> ")
	fmt.Scanf("%d", num)

	processes, err := monitor.GetProcesses((*tcp)[num].Inode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(*processes)
}
