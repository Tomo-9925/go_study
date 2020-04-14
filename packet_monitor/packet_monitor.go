package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/AkihiroSuda/go-netfilter-queue"
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

// キューの削除処理

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

	for {
		select {
		// パケットが届いたとき
		case p := <-packets:
			fmt.Println(p.Packet)
			p.SetVerdict(netfilter.NF_ACCEPT) // パケットを透過
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
