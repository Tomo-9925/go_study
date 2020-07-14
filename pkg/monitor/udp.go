package monitor

/*
ソースコードの構造や処理内容についてはtcp.goと類似しているため，処理内容の説明についてはtcp.goを参考にしてください．
共通化できれば嬉しいですが，関数名と引数を考える能がないので許してください…
*/

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"unsafe"

	"github.com/tomo-9925/go_study/pkg/utility"
)

var (
	udpFile string = ProcRoot + "/net/udp"
)

// UDPData は`/proc/net/tcp`の内容（一部のみ）
type UDPData struct {
	LocalIP    net.IP
	LocalPort  uint16
	RemoteIP   net.IP
	RemotePort uint16
	Inode      uint32
}

func (u UDPData) String() string {
	return fmt.Sprintf("{Local: %v:%d, Remote: %v:%d, inode: %d}", u.LocalIP, u.LocalPort, u.RemoteIP, u.RemotePort, u.Inode)
}

// GetAllUDPData は`/proc/net/udp`から取得した情報をUDPData構造体の入ったスライスで返却
func GetAllUDPData() ([]*UDPData, error) {
	f, err := os.Open(udpFile)
	var entries []*UDPData
	if err != nil {
		return entries, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return entries, err
	}

	s := strings.FieldsFunc(*(*string)(unsafe.Pointer(&b)), utility.Split)
	s = s[15:]
	for len(s) != 0 {
		localIP := utility.ParseIP(s[1])
		localPort := utility.ParsePort(s[2])
		remoteIP := utility.ParseIP(s[3])
		remotePort := utility.ParsePort(s[4])
		inode := utility.ParseInode(s[13])
		entry := UDPData{localIP, localPort, remoteIP, remotePort, inode}
		entries = append(entries, &entry)
		s = s[17:]
	}

	return entries, nil
}
