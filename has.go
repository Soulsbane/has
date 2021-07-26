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

func isMatch(dirName string, nameToSearchFor string, info fs.FileInfo) {
	if isValidPath(dirName) {
		if info.Name() == nameToSearchFor {
			if isSymbolicLink(info) {
				linkPath, err := filepath.EvalSymlinks(dirName)
				dirName = linkPath

				pathMatches[dirName] = linkPath

				if err != nil {
					fmt.Println(err)
				}
			}

			if !info.IsDir() {
				pathMatches[dirName] = ""
			}
		}
	}
}

func lookPath(fileName string) {
	stat, err := os.Lstat(fileName)

	if err == nil {

		if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
			linkPath, _ := filepath.EvalSymlinks(fileName)
			pathMatches[fileName] = linkPath
		} else {
			pathMatches[fileName] = ""
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
	var mutex = &sync.Mutex{}

	if !noPath {
		// FIXME: Need more tests to see which approach is faster
		/*path, _ := exec.LookPath(nameToSearchFor)
		lookPath(path)*/
		envPath := getAdditionalPaths()
		searchPaths = append(searchPaths, envPath...)
	}

	for _, dirToSearch := range searchPaths {
		walkFn := func(pathname string, fi os.FileInfo) error {
			mutex.Lock()
			isMatch(pathname, nameToSearchFor, fi)
			mutex.Unlock()

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
