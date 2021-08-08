package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
	"github.com/saracen/walker"
)

var searchPaths = []string{
	"/usr/bin",
	"/usr/sbin",

	"/usr/local/bin",
	"/usr/local/sbin",

	"/bin",
	"/sbin",
	"/opt/bin",
	"/usr/share", // Needs permission
}

var pathMatches = map[string]string{}

func isValidPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func isSymbolicLink(info fs.FileInfo) bool {
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
}

func isFileExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}

func addMatches(dirName string, nameToSearchFor string, info fs.FileInfo) {
	var mutex = &sync.Mutex{}

	if isValidPath(dirName) { // FIXME: Cleanup multiple calls
		if info.Name() == nameToSearchFor {
			if isSymbolicLink(info) {
				linkPath, err := filepath.EvalSymlinks(dirName)

				if isFileExecutable(info.Mode()) {
					mutex.Lock()
					pathMatches[dirName] = linkPath
					mutex.Unlock()

					if err != nil {
						fmt.Println(err)
					}
				}
			} else {
				if !info.IsDir() && isFileExecutable(info.Mode()) {
					mutex.Lock()
					pathMatches[dirName] = ""
					mutex.Unlock()
				}
			}
		}
	}
}

func getAdditionalPaths() []string {
	path, variableExists := os.LookupEnv("PATH")
	var paths []string

	if variableExists {
		paths = strings.Split(path, ":")
	}

	return paths
}

func findExecutable(nameToSearchFor string, noPath bool) {
	if !noPath {
		envPath := getAdditionalPaths()
		searchPaths = append(searchPaths, envPath...)
	}

	for _, dirToSearch := range searchPaths {
		walkFn := func(path string, fileInfo os.FileInfo) error {
			addMatches(path, nameToSearchFor, fileInfo)
			return nil
		}

		errorCallbackOption := walker.WithErrorCallback(func(pathname string, err error) error {
			if os.IsPermission(err) {
				return nil // INFO: Ignore permission errors
			}

			return err // INFO: Stop on all other errors
		})

		walker.Walk(dirToSearch, walkFn, errorCallbackOption)
	}
}

func colorizePath(name string, ugly bool) string {
	if ugly {
		return name
	}

	dir := color.BlueString(filepath.Dir(name))
	base := color.GreenString(filepath.Base(name))

	return path.Join(dir, base)
}

func listMatches(ugly bool) {
	if len(pathMatches) > 0 {
		for path, linkPath := range pathMatches {
			if len(linkPath) > 0 {
				fmt.Printf("%s => %s\n", colorizePath(path, ugly), colorizePath(linkPath, ugly))
			} else {
				fmt.Println(colorizePath(path, ugly))
			}
		}
	} else {
		fmt.Println("No files found!")
	}
}

func main() {
	var args struct {
		FileName string `arg:"positional, required"`
		NoPath   bool   `arg:"-n, --no-path" default:"false" help:"Include directories in user's $PATH."`
		Ugly     bool   `arg:"-u, --ugly" default:"false" help:"Remove colorized output. Yes it's ugly."`
	}

	arg.MustParse(&args)
	findExecutable(args.FileName, args.NoPath)
	listMatches(args.Ugly)
}
