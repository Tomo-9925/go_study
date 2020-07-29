package monitor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"unsafe"

	gnq "github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	localIP []net.IP
)

// Socket はprocファイルシステムのnetから取得できるデータの構造体です．
type Socket struct {
	LayerType             gopacket.LayerType
	LocalIP, RemoteIP     net.IP
	LocalPort, RemotePort uint16
	Inode                 uint32
}

func (s *Socket) String() string {
	return fmt.Sprintf("{LayerType: %s, Src: %v:%d, Dst: %v:%d, inode: %d}", s.LayerType.String(), s.LocalIP, s.LocalPort, s.RemoteIP, s.RemotePort, s.Inode)
}

// GetSocket はパケットの情報から`/proc/net`の情報を取得してSocket構造体を作成します．
func GetSocket(p gnq.NFPacket) (*Socket, error) {
	var file string
	var s Socket
	var parseEntry func(string) Socket
	var port = map[string]uint16{"srcPort": 0, "dstPort": 0}

	// パケットからトランスポート層のプロトコル情報を取得
	t := p.Packet.TransportLayer()
	// トランスポート層が無い（ICMPなど）ときは無視
	if t == nil {
		// TODO: icmpパケットの処理方法を書く
		return &s, errors.New("Protocol not supported")
	}
	s.LayerType = t.LayerType()
	// トランスポート層のプロトコルによって読み込むファイル，パース方法を変更し，ポート番号を取得
	if s.LayerType == layers.LayerTypeTCP {
		file = tcpFile
		parseEntry = parseTCPEntry
		tcp, _ := p.Packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
		port["srcPort"] = uint16(tcp.SrcPort)
		port["dstPort"] = uint16(tcp.DstPort)
	} else if s.LayerType == layers.LayerTypeUDP {
		file = udpFile
		parseEntry = parseUDPEntry
		udp, _ := p.Packet.Layer(layers.LayerTypeUDP).(*layers.UDP)
		port["srcPort"] = uint16(udp.SrcPort)
		port["dstPort"] = uint16(udp.DstPort)
	} else {
		return &s, errors.New("Protocol not supported")
	}

	// パケットからSocket構造体の項目を埋める
	if len(localIP) == 0 {
		getLocalIP()
	}
	ip, _ := p.Packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
	for _, l := range localIP {
		if ip.SrcIP.Equal(l) {
			s.LocalPort = port["srcPort"]
			s.RemoteIP = ip.DstIP
			s.RemotePort = port["dstPort"]
		} else if ip.DstIP.Equal(l) {
			s.LocalPort = port["dstPort"]
			s.RemoteIP = ip.SrcIP
			s.RemotePort = port["srcPort"]
		}
	}
	if s.LocalPort == 0 {
		return &s, errors.New("Local IP address and local port not found")
	}

	// ファイルの読み込み
	f, err := os.Open(file)
	if err != nil {
		return &s, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f) // ファイルをすべてを読み込む
	if err != nil {
		return &s, err
	}

	// エントリの検索処理
	rawEntries := strings.Split(*(*string)(unsafe.Pointer(&b)), "\n") // 改行でスライスを作成
	rawEntries = rawEntries[1 : len(rawEntries)-1]                    // いらない行をスキップ
	for i, rawEntry := range rawEntries {
		entry := parseEntry(rawEntry)
		fmt.Printf("Entry%d: %v\n", i, entry)
		if entry.LocalPort == s.LocalPort && entry.RemoteIP.Equal(s.RemoteIP) && entry.RemotePort == s.RemotePort {
			s = entry
			break
		} else if entry.LocalPort == s.LocalPort || entry.LocalPort == s.RemotePort || entry.RemotePort == s.LocalPort || entry.RemotePort == s.RemotePort {
			s = entry
		}
	}
	if s.Inode == 0 {
		return &s, errors.New("inode not found")
	}

	return &s, nil
}

func getLocalIP() error {
	localIP = nil // スライスの初期化

	// コンテナからLocalIPを取得する
	if len(containerInformations) == 0 {
		storeContainerInfo()
	}
	for _, c := range containerInformations {
		if c.IP != nil {
			localIP = append(localIP, c.IP)
		}
	}

	// インタフェースからLocalIPを取得する
	ifaces, _ := net.Interfaces() // インタフェース情報の取得
	for _, i := range ifaces {
		addrs, err := i.Addrs() // インタフェースのIPアドレスを取得
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// IPv4アドレスだけ取得
			if ip.To4() != nil {
				localIP = append(localIP, ip)
			}
		}
	}
	localIP = append(localIP, net.ParseIP("127.0.0.53"))

	return nil
}
