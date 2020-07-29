package monitor

import (
	"io/ioutil"
	"os"
	"strings"
	"unsafe"

	"github.com/google/gopacket/layers"
	"github.com/tomo-9925/go_study/pkg/utility"
)

var (
	tcpFile string = ProcRoot + "/net/tcp"
	// tcpFile string = "./tcp" // 検証用
)

// GetAllTCPData は`/proc/net/tcp`から取得した情報をSocket構造体の入ったスライスで返却する
func GetAllTCPData() ([]*Socket, error) {
	// ファイルの読み込み
	f, err := os.Open(tcpFile)
	var entries []*Socket
	if err != nil {
		return entries, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f) // ファイルをすべてを読み込む
	if err != nil {
		return entries, err
	}

	// データの格納処理
	rawEntries := strings.Split(*(*string)(unsafe.Pointer(&b)), "\n") // 改行でスライスを作成
	rawEntries = rawEntries[1 : len(rawEntries)-1]                    // いらない行をスキップ
	for _, rawEntry := range rawEntries {
		entry := parseTCPEntry(rawEntry)
		entries = append(entries, &entry)
	}

	return entries, nil
}

func parseTCPEntry(e string) Socket {
	s := strings.FieldsFunc(e, utility.Split) // " "と":"で文字列分割
	localIP := utility.ParseIP(s[1])
	localPort := utility.ParsePort(s[2])
	remoteIP := utility.ParseIP(s[3])
	remotePort := utility.ParsePort(s[4])
	inode := utility.ParseInode(s[13])
	return Socket{layers.LayerTypeTCP, localIP, remoteIP, localPort, remotePort, inode}
}
