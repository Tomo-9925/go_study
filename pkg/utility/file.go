package utility

import (
	"os"
	"path/filepath"
)

// GetFilePath はファイルの有無の確認をします
func GetFilePath(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return path, err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return path, err
	}
	return GetPathFromLink(path), err
}

// GetPathFromLink はシンボリックリンクから実際のファイルパスを返します
func GetPathFromLink(link string) string {
	path, err := os.Readlink(link)
	if err != nil {
		return link
	}
	return GetPathFromLink(path)
}
