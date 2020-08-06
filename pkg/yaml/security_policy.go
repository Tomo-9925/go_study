package yaml

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"strings"

	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
	"github.com/tomo-9925/go_study/pkg/monitor"
	"gopkg.in/yaml.v3"
)

// Network はコンテナのネットワークに関する必要な情報
// type Network struct {
// 	ID     string
// 	Driver string
// 	Subnet net.IPNet
// }

// Filter はフィルターに関する必要な情報
type Filter struct {
	Processes []*Process
	Sockets   []*monitor.Socket
}

// Firewall はコンテナのファイアウォールの設定情報
type Firewall struct {
	Type      string
	Container *monitor.Container
	Filters   []*Filter
}

// Process はコンテナのプロセスに関する必要な情報
type Process struct {
	Executable string
	Path       string
	Pid        int
	Inode      uint64
}

// ParseSecurityPolicy はYAMLで書かれたセキュリティポリシーを適切な構造体にパースします．
func ParseSecurityPolicy(path string) ([]*Firewall, error) {
	var firewalls []*Firewall
	var err error

	// 指定されたPathを開き，全体をmap[string]interface{}としてパース
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return firewalls, err
	}
	sp := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &sp)
	if err != nil {
		return firewalls, err
	}

	// Firewallsの情報を取得
	if _, exist := sp["firewalls"]; !exist {
		err = errors.New("Firewallsが見つかりませんでした．")
	} else if spFirewalls, ok := sp["firewalls"].([]interface{}); !ok {
		err = errors.New("Firewallsを配列として展開することができませんでした．")
	} else {
		for _, spFirewallInterface := range spFirewalls {
			spFirewall, ok := spFirewallInterface.(map[string]interface{})
			if !ok {
				err = errors.New("Firewallsの配列内が連想配列ではありません．")
				break
			}
			firewall, err := parseFirewall(spFirewall)
			if err != nil {
				return firewalls, err
			}
			firewalls = append(firewalls, firewall)
			logrus.Debugln("Firewalls:", firewalls)
		}
	}

	return firewalls, err
}

func parseFirewall(spFirewall map[string]interface{}) (*Firewall, error) {
	var firewall Firewall

	// Typeの取得
	if _, exist := spFirewall["type"]; !exist {
		firewall.Type = "allowlist" // デフォルトはallowlsitとする．
	} else if spType, ok := spFirewall["type"].(string); !ok {
		return &firewall, errors.New("FirewallのTypeは文字列で記入してください．")
	} else if strings.EqualFold(spType, "allowlist") {
		firewall.Type = "allowlist"
	} else if strings.EqualFold(spType, "denylist") {
		firewall.Type = "denylist"
	} else {
		return &firewall, errors.New("FirewallのTypeにallowlist，denylist以外の文字列が記入されています．")
	}
	logrus.Debugln("Type:", firewall.Type)

	// コンテナ情報を取得
	if _, exist := spFirewall["container"]; !exist {
		return &firewall, errors.New("FirewallのContainerの情報が記入されていません．")
	}
	spContainer, ok := spFirewall["container"].(map[string]interface{})
	if !ok {
		return &firewall, errors.New("FirewallのContainerを連想配列として展開することができませんでした．")
	}
	container, err := parseContainer(spContainer)
	if err != nil {
		return &firewall, err
	}
	firewall.Container = container
	logrus.Debugln("Container:", firewall.Container)

	// フィルタ情報を取得
	if _, exist := spFirewall["filters"]; !exist {
		return &firewall, errors.New("FirewallのFilterの情報が記入されていません．")
	} else if spFilters, ok := spFirewall["filters"].([]interface{}); !ok {
		return &firewall, errors.New("FirewallのFiltersを配列として展開することができませんでした．")
	} else {
		var filters []*Filter
		for _, spFilterInterface := range spFilters {
			var filter *Filter
			spFilter, ok := spFilterInterface.(map[string]interface{})
			if !ok {
				return &firewall, errors.New("FirewallのFiltersの配列内が連想配列ではありません．")
			}
			filter, err = parseFilter(spFilter)
			if err != nil {
				return &firewall, err
			}
			filters = append(filters, filter)
			logrus.Debugln("Filters:", filters)
		}
		firewall.Filters = filters
	}
	logrus.Debugln("Firewall:", firewall)

	return &firewall, nil
}

func parseContainer(spContainer map[string]interface{}) (*monitor.Container, error) {
	var container monitor.Container

	if _, exist := spContainer["name"]; !exist {
		if _, exist := spContainer["id"]; !exist {
			return &container, errors.New("FirewallのContainerにIDまたはNameが記入されていません")
		}
		container.ID = spContainer["id"].(string)
		logrus.Debugln("ID:", container.ID)
	} else {
		container.Name = spContainer["name"].(string)
		logrus.Debugln("Name:", container.Name)
	}

	// 取得したコンテナ情報から補完

	return &container, nil
}

func parseFilter(spFilter map[string]interface{}) (*Filter, error) {
	var filter Filter

	// プロセス情報の取得
	if _, exist := spFilter["processes"]; !exist {
		return &filter, errors.New("FirewallのFilterのProcessesが記入されていません．")
	} else if spProcesses, ok := spFilter["processes"].([]interface{}); !ok {
		return &filter, errors.New("FirewallのFilterのProcessesを配列として展開することができませんでした．")
	} else {
		var processes []*Process
		for _, spProcessInterface := range spProcesses {
			spProcess, ok := spProcessInterface.(map[string]interface{})
			if !ok {
				return &filter, errors.New("FirewallのFiltersの配列内が連想配列ではありません．")
			}
			process, err := parseProcess(spProcess)
			if err != nil {
				return &filter, err
			}
			processes = append(processes, process)
			logrus.Debugln("Processes:", processes)
		}
	}

	// ソケット情報の取得
	if _, exist := spFilter["sockets"]; !exist {
		return &filter, errors.New("FirewallのFilterのSocketsが記入されていません．")
	} else if Sockets, ok := spFilter["sockets"].([]interface{}); !ok {
		return &filter, errors.New("FirewallのFilterのSocketsを配列として展開することができませんでした．")
	} else {
		var sockets []*monitor.Socket
		for _, spSocketInterface := range Sockets {
			spSocket, ok := spSocketInterface.(map[string]interface{})
			if !ok {
				return &filter, errors.New("FirewallのFiltersの配列内が連想配列ではありません．")
			}
			socket, err := parseSocket(spSocket)
			if err != nil {
				return &filter, err
			}
			sockets = append(sockets, socket)
			logrus.Debugln("Sockets:", sockets)
		}
	}

	return &filter, nil
}

func parseProcess(spProcess map[string]interface{}) (*Process, error) {
	var process Process

	if _, exist := spProcess["path"]; !exist {
		if _, exist := spProcess["executable"]; !exist {
			return &process, errors.New("FirewallのFilterのProcessにPathまたはExecutableが記入されていません．")
		} else if spExecutable, ok := spProcess["executable"].(string); !ok {
			return &process, errors.New("FirewallのFilterのProcessのExecutableが文字列ではありません．")
		} else {
			process.Executable = spExecutable
			logrus.Debugln("Executable:", process.Executable)
		}
	} else {
		psPath, ok := spProcess["path"].(string)
		if !ok {
			return &process, errors.New("FirewallのFilterのProcessのPathが文字列ではありません．")
		}
		process.Path = psPath
		logrus.Debugln("Path:", process.Path)
	}

	return &process, nil
}

func parseSocket(spSocket map[string]interface{}) (*monitor.Socket, error) {
	var socket monitor.Socket

	// プロトコル情報の取得
	if _, exist := spSocket["protocol"]; !exist {
		return &socket, errors.New("FirewallのFilterのSocketにProtocolが記入されていません．")
	} else if spProtocol, ok := spSocket["protocol"].(string); !ok {
		return &socket, errors.New("FirewallのFilterのSocketのProtocolが文字列ではありません．")
	} else if strings.EqualFold(spProtocol, "ICMP") {
		socket.LayerType = layers.LayerTypeICMPv4
	} else if strings.EqualFold(spProtocol, "TCP") {
		socket.LayerType = layers.LayerTypeTCP
	} else if strings.EqualFold(spProtocol, "UDP") {
		socket.LayerType = layers.LayerTypeUDP
	} else {
		return &socket, errors.New("FirewallのFilterのSocketのProtocolに未対応のプロトコルが指定されています．")
	}
	logrus.Debugln("Type:", socket.LayerType)

	// ローカルIPアドレス情報を取得（いらないかな…？）

	// リモートIPアドレス情報を取得
	if _, exist := spSocket["remote_container_name"]; exist {
		logrus.Debugln("RemoteIP: specified container name")
	} else if _, exist := spSocket["remote_ip"]; exist {
		if spRemoteIP, ok := spSocket["remote_ip"].(string); !ok {
			return &socket, errors.New("FirewallのFilterのSocketのremote_ipが文字列ではありません．")
		} else if socket.RemoteIP = net.ParseIP(spRemoteIP); socket.RemoteIP == nil {
			return &socket, errors.New("FirewallのFilterのSocketのremote_ipが適切な文字列ではありません．")
		}
		logrus.Debugln("RemoteIP:", socket.RemoteIP)
	} else {
		logrus.Debugln("remote_container_nameとremote_ipは見つかりませんでした．")
	}

	// ローカルとリモートのポート情報を取得
	if socket.LayerType != layers.LayerTypeICMPv4 {
		if _, exist := spSocket["local_port"]; exist {
			spLocalPort, ok := spSocket["local_port"].(int)
			if !ok {
				return &socket, errors.New("FirewallのFilterのSocketのlocal_ipが適切な数値ではありません．")
			}
			socket.LocalPort = uint16(spLocalPort)
			logrus.Debugln("LocalPort:", socket.LocalPort)
		}
		if _, exist := spSocket["remote_port"]; exist {
			spRemotePort, ok := spSocket["remote_port"].(int)
			if !ok {
				return &socket, errors.New("FirewallのFilterのSocketのremote_ipが適切な数値ではありません．")
			}
			socket.RemotePort = uint16(spRemotePort)
			logrus.Debugln("RemotePort:", socket.RemotePort)
		}
	}

	return &socket, nil
}

// デバッグ用インタフェース確認関数
func check(object interface{}) {
	fmt.Println(reflect.TypeOf(object))
	fmt.Println(object)
}
