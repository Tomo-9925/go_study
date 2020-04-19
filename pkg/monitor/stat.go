package monitor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

// ProcessWithPath は通信を遮断する際に必要となるプロセス情報を格納した構造体です．
// go-psのUnixProcess構造体（pid, ppid, state, pgrp, sid, binary）の情報に加えて，実行ファイルのパスの情報を追加しています．
type ProcessWithPath struct {
	ps.Process
	path string
}

func (p ProcessWithPath) String() string {
	return fmt.Sprintf("{Pid: %d, PPid: %d, Executable: %s, Path: %s}", p.Pid(), p.PPid(), p.Executable(), p.path)
}

// GetProcess は引数のinodeの数値をもとにプロセスを検索し，ProcessWithPath型のスライスでで情報を返却します．
func GetProcess(inode uint32) ([]*ProcessWithPath, error) {
	var processInfo []*ProcessWithPath

	// すべてのプロセス情報を取得
	process, err := ps.Processes()
	if err != nil {
		fmt.Println(err)
		return processInfo, err
	}

	// すべてのプロセス情報から指定されたinodeがあるか調査
	for _, process := range process {
		fdPath := ProcRoot + "/" + strconv.Itoa(process.Pid()) + "/fd"
		if ExistInode(fdPath, inode) {
			path := GetProcessPath(process.Pid())
			p := ProcessWithPath{process, path}
			processInfo = append(processInfo, &p)
		}
	}
	if len(processInfo) == 0 {
		err = errors.New("couldn't find process")
		fmt.Println(err)
		return processInfo, err
	}

	return processInfo, err
}

// ExistInode は指定されたプロセスのfdディレクトリ内に指定されたinodeを情報を含むシンボリックリンクの有無を確認します．
func ExistInode(fdPath string, inode uint32) bool {
	// ファイルディスクリプタの情報を取得
	files, err := ioutil.ReadDir(fdPath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// inodeの数値を文字列型に変換
	i := strconv.FormatUint(uint64(inode), 10)

	// ファイルディスクリプタごとに指定されたinodeがあるかを確認
	for _, file := range files {
		fd := filepath.Join(fdPath, file.Name())
		str, err := os.Readlink(fd)
		if err != nil {
			fmt.Println(err)
		} else if strings.Contains(str, i) {
			return true
		}
	}
	return false
}

// GetProcessPath は指定されたプロセスIDからプロセスのフルパスを返します．
func GetProcessPath(pid int) string {
	exe := ProcRoot + "/" + strconv.Itoa(pid) + "/exe"
	str, err := os.Readlink(exe)
	if err != nil {
		fmt.Println(err)
	}
	return str
}
