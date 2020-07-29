package monitor

/*
ソースコードの構造や処理内容についてはtcp.goと類似しているため，処理内容の説明についてはtcp.goを参考にしてください．
*/

import (
	"io/ioutil"
	"os"
	"strings"
	"unsafe"

	"github.com/google/gopacket/layers"
	"github.com/tomo-9925/go_study/pkg/utility"
)

var (
	udpFile string = ProcRoot + "/net/udp"
)

// GetAllUDPData は`/proc/net/udp`から取得した情報をUDPData構造体の入ったスライスで返却
func GetAllUDPData() ([]*Socket, error) {
	f, err := os.Open(udpFile)
	var entries []*Socket
	if err != nil {
		return entries, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return entries, err
	}

	// データの格納処理
	rawEntries := strings.Split(*(*string)(unsafe.Pointer(&b)), "\n") // 改行でスライスを作成
	rawEntries = rawEntries[1 : len(rawEntries)-1]                    // いらない行をスキップ
	for _, rawEntry := range rawEntries {
		entry := parseUDPEntry(rawEntry)
		entries = append(entries, &entry)
	}

	return entries, nil
}

func parseUDPEntry(e string) Socket {
	s := strings.FieldsFunc(e, utility.Split) // " "と":"で文字列分割
	localIP := utility.ParseIP(s[1])
	localPort := utility.ParsePort(s[2])
	remoteIP := utility.ParseIP(s[3])
	remotePort := utility.ParsePort(s[4])
	inode := utility.ParseInode(s[13])
	return Socket{layers.LayerTypeUDP, localIP, remoteIP, localPort, remotePort, inode}
}
