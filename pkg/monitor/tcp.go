package monitor

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
	tcpFile               string = ProcRoot + "/net/tcp"
	updateTCPDataInterval uint   = 1000
	// tcpFile string = "./tcp" // 検証用
)

// TCPData は`/proc/net/tcp`の内容（全部やるのは辛いので一部だけ）
type TCPData struct {
	EntryNum   uint16
	LocalIP    net.IP
	LocalPort  uint16
	RemoteIP   net.IP
	RemotePort uint16
	// ConnectionState
	// TransmitQueue
	// ReceiveQueue
	// TimerActive
	// NumberOfJiffiesUntilTimeExpires
	// NumverOfUnrecoveredRTOTimeouts
	// UID
	// Unanswered0WindowProves
	Inode uint32
	// SocketReferenceCount
	// LocationOfSocketInMemory
	// RetransmiTimeout
	// PredictedTickOfSoftClock
	// AckQuick
	// SendingCongestionWindow
	// SlowStartSizeThreshold
}

func (t TCPData) String() string {
	return fmt.Sprintf("{EntryNum: %d, Local: %v:%d, Remote: %v:%d, inode: %d}", t.EntryNum, t.LocalIP, t.LocalPort, t.RemoteIP, t.RemotePort, t.Inode)
}

// GetAllTCPData は`/proc/net/tcp`から取得した情報をTCPData構造体の入ったスライスで返却
func GetAllTCPData() ([]*TCPData, error) {
	// ファイルの読み込み
	f, err := os.Open(tcpFile)
	var entries []*TCPData
	if err != nil {
		return entries, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f) // ファイルをすべてを読み込む
	if err != nil {
		return entries, err
	}

	// データの格納処理
	s := strings.FieldsFunc(*(*string)(unsafe.Pointer(&b)), utility.Split) // " "と":"，"\n"で文字列分割
	s = s[12:]                                                             // インデックス行の削除
	for len(s) != 0 {
		entryNum := utility.ParseEntryNum(s[0])
		localIP := utility.ParseIP(s[1])
		localPort := utility.ParsePort(s[2])
		remoteIP := utility.ParseIP(s[3])
		remotePort := utility.ParsePort(s[4])
		inode := utility.ParseInode(s[13])
		entry := TCPData{entryNum, localIP, localPort, remoteIP, remotePort, inode}
		entries = append(entries, &entry)
		s = s[21:] // スライスの頭を次の1行に移動
	}

	return entries, nil
}
