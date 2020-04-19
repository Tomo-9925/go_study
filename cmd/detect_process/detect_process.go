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

	for i, d := range tcp {
		fmt.Println(i, d)
	}
	var num uint16
	fmt.Print("どのTCPの情報を取得しますか．インデックスを指定してください．\n> ")
	fmt.Scanf("%d", num)

	process, err := monitor.GetProcess(tcp[num].Inode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(process)
}
