package utility

import (
	gnq "github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket/layers"
	"os"
	"bufio"
	"strings"
	"strconv"
	"unsafe"
	"github.com/mitchellh/go-ps"
	"fmt"
)

var (
	tcpFile string = "/proc/net/tcp"
	udpFile string = "/proc/net/udp"
)

//CheckProtocol returns TCP:1 UDP:2 Others:0
func CheckProtocol(p gnq.NFPacket) uint16{
	if tcpLayer := p.Packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		return 1;
	}else if udpLayer := p.Packet.Layer(layers.LayerTypeUDP); udpLayer !=nil {
		return 2;
	}else{
		return 0;
	}
}	

//CheckSrcPort returns source port 
func CheckSrcPort(p gnq.NFPacket,protocolNum uint16) uint16 {
	var srcPort uint16
	if protocolNum == 1 {
		tcpLayer := p.Packet.Layer(layers.LayerTypeTCP)
		// Get actual TCP data from this layer
		tcp, _ := tcpLayer.(*layers.TCP)
		srcPort = (uint16)(tcp.SrcPort)
	}else if protocolNum == 2 {
		udpLayer := p.Packet.Layer(layers.LayerTypeUDP);
		// Get actual UDP data from this layer
		udp,_ := udpLayer.(*layers.UDP)
		srcPort = (uint16)(udp.SrcPort)
	}	
	return srcPort
}

//GetInode returns inode
func GetInode(protocolNum uint16,srcPort uint16)(uint32,error){
		var filename string
		myPid :=os.Getpid()
		//protocolによって読み込みファイルを変える TCPなら /proc/{pid}/net/tcp , UDPなら /proc/{pid}/net/udp
		if protocolNum == 1{
			filename = "tcp"
		}else if protocolNum == 2{
			filename = "udp"
		}
	
		//全てのプロセスを取得する
		processes,err := ps.Processes()

		if err != nil{
			return 0,err
		}

		//全てのプロセスの/proc/pid/net/tcp(udp)を見ていく
		for _,process := range processes{
			if process.Pid() == myPid {
				continue
			}
			filePath :="/proc"+"/"+strconv.Itoa(process.Pid())+"/net"+"/"+filename
			file,err := os.Open(filePath)
			if err != nil{
				return 0,err
			}

			//終了時にファイルをクローズする
			defer file.Close()
			//ioutilだと一括読み込みになるのでbufioを使用
			scanner := bufio.NewScanner(file)

			//1行目は不要なためスキップ
			scanner.Scan()

			//１行ずつ読み込む
			for scanner.Scan(){
				line := scanner.Text()

				// データの格納処理
				str := strings.FieldsFunc(*(*string)(unsafe.Pointer(&line)), Split) // " "と":"，"\n"で文字列分割
				localPort := ParsePort(str[2])

				//port番号が取得したいものと一致すればinodeを取得し、返す
				if localPort == srcPort{
					fmt.Println(process.Pid())
					inode := ParseInode(str[13])
					return inode,nil
				}
		}
	
		if err := scanner.Err(); err!=nil{
			return 0,err
		}
	}
		return 0,nil
}