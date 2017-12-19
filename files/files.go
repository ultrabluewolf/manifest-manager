package files

import (
	"os"
)

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func IsDir(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return true
	}
	return fileInfo.IsDir()
}
