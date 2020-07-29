package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket/layers"
	"github.com/tomo-9925/go_study/pkg/monitor"
)

// 定数
const (
	Debug             = true // デバッグ機能
	QueueNum          = 2    // キュー番号
	MaxPacketsInQueue = 100  // キューのサイズ
)

// デバッグ用プリント関数
func printD(msg string) {
	if Debug {
		fmt.Println("debug:", msg)
	}
}

func main() {
	printD("debug mode")

	// 変数の宣言
	var err error                                                                                   // エラー
	nfq, err := netfilter.NewNFQueue(QueueNum, MaxPacketsInQueue, netfilter.NF_DEFAULT_PACKET_SIZE) // キューの定義
	c := make(chan os.Signal, 1)                                                                    // sigintを待ち受けるchan

	// キューの作成失敗
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 終了時の処理
	defer nfq.Close() // キューの監視終了
	printD("create and bind to queue specified by QueueNum")

	// SIGINTのイベントをフック
	signal.Notify(c, os.Interrupt)
	printD("signal handling")

	// パケットをキューから取得するchanの作成
	packets := nfq.GetPackets()
	printD("packet handling")

	count := 0

	for {
		select {
		// パケットが届いたとき
		case p := <-packets:
			count++
			fmt.Printf("%d番目のパケット\n", count)
			fmt.Printf("パケットの概要:\n%v\n", p.Packet)
			fmt.Println("パケットの解析結果:")
			if p.Packet.NetworkLayer().LayerType() == layers.LayerTypeIPv4 {
				s, err := monitor.GetSocket(p)
				if err != nil && Debug {
					fmt.Println(err)
				}
				fmt.Printf("Socket: %+v\n", s)
				process, err := monitor.GetProcess(s.Inode)
				if err != nil && Debug {
					fmt.Println(err)
				}
				fmt.Printf("Process: %+v\n", process)
			} else {
				printD("IPv4以外のネットワークレイヤーのプロトコルを使用した通信を観測しました．")
			}
			p.SetVerdict(netfilter.NF_ACCEPT) // パケットを透過
			fmt.Printf("\n\n")
		// SIGINTを検知したとき
		case sig := <-c:
			fmt.Println(sig)
			close(c)
			nfq.Close()
			printD("nfq unbinded")
			os.Exit(130)
		}
	}
}
