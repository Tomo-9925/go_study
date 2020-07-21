package utility

import (
	"bufio"
	"encoding/json"
	gnq "github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket/layers"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

var (
	tcpFile string = "/proc/net/tcp"
	udpFile string = "/proc/net/udp"
)

type Config struct {
	State struct {
		Pid  int    `json:Pid`
		ID   int    `json:ID`
		Name string `json:Name`
	} `json:State`
}

//CheckProtocol returns TCP:1 UDP:2 Others:0
func CheckProtocol(p gnq.NFPacket) uint16 {
	if tcpLayer := p.Packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		return 1
	} else if udpLayer := p.Packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		return 2
	} else {
		return 0
	}
}

//CheckSrcPort returns source port
func CheckSrcPort(p gnq.NFPacket, protocolNum uint16) uint16 {
	var srcPort uint16
	if protocolNum == 1 {
		tcpLayer := p.Packet.Layer(layers.LayerTypeTCP)
		// Get actual TCP data from this layer
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort = (uint16)(tcp.SrcPort)
	} else if protocolNum == 2 {
		udpLayer := p.Packet.Layer(layers.LayerTypeUDP)
		// Get actual UDP data from this layer
		udp, _ := udpLayer.(*layers.UDP)
		srcPort = (uint16)(udp.SrcPort)
	}
	return srcPort
}

//GetInode returns inode
func GetInode(protocolNum uint16, srcPort uint16) (uint32, error) {
	var filename string
	//protocolによって読み込みファイルを変える TCPなら /proc/{pid}/net/tcp , UDPなら /proc/{pid}/net/udp
	if protocolNum == 1 {
		filename = "tcp"
	} else if protocolNum == 2 {
		filename = "udp"
	}

	//dockerのコンテナ情報が記載されているディレクトリを開く
	configDirPath := "/var/lib/docker/containers"
	containerFiles, err := ioutil.ReadDir(configDirPath)

	if err != nil {
		return 0, err
	}
	for _, containerFile := range containerFiles {
		bytes, err := ioutil.ReadFile(configDirPath + "/" + containerFile.Name() + "/config.v2.json")
		if err != nil {
			return 0, err
		}

		//config構造体へ必要な情報を当てはめる
		var config Config
		if err := json.Unmarshal(bytes, &config); err != nil {
			return 0, err
		}

		//containerが起動していない(Pidが0)なら次のファイルへ
		if config.State.Pid == 0 {
			continue
		}
		//dockerコンテナで使用しているPid以下のnet/tcp(udp)を見る
		filePath := "/proc" + "/" + strconv.Itoa(config.State.Pid) + "/net" + "/" + filename

		inode, err := func() (uint32, error) {
			file, err := os.Open(filePath)
			if err != nil {
				return 0, err
			}
			//終了時にファイルをクローズする
			defer file.Close()
			//ioutilだと一括読み込みになるのでbufioを使用
			scanner := bufio.NewScanner(file)

			//1行目は不要なためスキップ
			scanner.Scan()

			//１行ずつ読み込む
			for scanner.Scan() {
				line := scanner.Text()

				// データの格納処理
				str := strings.FieldsFunc(*(*string)(unsafe.Pointer(&line)), Split) // " "と":"，"\n"で文字列分割
				localPort := ParsePort(str[2])

				//port番号が取得したいものと一致すればinodeを取得し、返す
				if localPort == srcPort {
					inode := ParseInode(str[13])
					return inode, nil
				}
			}
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, nil
		}()

		if err != nil {
			return 0, err
		}
		if inode != 0 {
			return inode, err
		}

	}
	return 0, nil
}
