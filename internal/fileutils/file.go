package fileutils

import (
	"os"
	"runtime"

	"github.com/duke-git/lancet/v2/fileutil"
)

func IsFileHidden(info os.DirEntry) bool {
	if runtime.GOOS != "windows" {
		return info.Name()[0:1] == "."
	}

	return false
}

// GetLinkPath returns the path of the link and a boolean indicating if the link destination path exists
func GetLinkPath(name string) (string, bool) {
	realPath, err := os.Readlink(name)

	if err != nil {
		return "", fileutil.IsExist(realPath)
	}

	return realPath, fileutil.IsExist(realPath)
}

func IsFileExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}
