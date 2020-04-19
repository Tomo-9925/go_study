package utility

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// Split は文字が" "または":"，"\n"であるかを判定します．
func Split(r rune) bool {
	return r == ' ' || r == ':' || r == '\n'
}

// RemoveZeroPadding は文字列の頭の"0"を取り除きます．
func RemoveZeroPadding(str string) {
	for len(str) > 1 {
		if str[0] != '0' {
			break
		}
		str = str[1:]
	}
}

// ParseEntryNum はprocファイルシステムから取得したエントリー番号をuint8型に変換します．
func ParseEntryNum(str string) uint16 {
	RemoveZeroPadding(str)
	s, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		fmt.Println(err)
	}
	return uint16(s)
}

// ParseIP はprocファイルシステムから取得したIPv4アドレスをnet.IP型に変換します．
func ParseIP(str string) net.IP {
	ip := make(net.IP, 4)
	RemoveZeroPadding(str)
	s, err := strconv.ParseUint(str, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	binary.LittleEndian.PutUint32(ip, uint32(s))
	return ip
}

// ParsePort はprocファイルシステムから取得したポート番号をuint16に変換します．
func ParsePort(str string) uint16 {
	RemoveZeroPadding(str)
	s, err := strconv.ParseUint(str, 16, 16)
	if err != nil {
		fmt.Println(err)
	}
	return uint16(s)
}

// ParseInode はprocファイルシステムから取得したinode番号をuint32に変換します．
func ParseInode(str string) uint32 {
	RemoveZeroPadding(str)
	s, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	return uint32(s)
}

/*
ParseXX系の関数をうまくまとめれれば嬉しい．けど，そこまで考えようという気が起きない…
*/
