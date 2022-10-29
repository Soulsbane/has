package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/Soulsbane/has/internal/fileutils"
	"github.com/Soulsbane/has/internal/paths"
	"github.com/alexflint/go-arg"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/fatih/color"
	"github.com/saracen/walker"
)

var pathMatches = map[string]string{}

func addMatches(dirName string, nameToSearchFor string, info fs.FileInfo) {
	var mutex = &sync.Mutex{}

	if fileutil.IsExist(dirName) {
		if info.Name() == nameToSearchFor && fileutils.IsFileExecutable(info.Mode()) && !info.IsDir() {
			if fileutil.IsLink(dirName) {
				linkPath, err := filepath.EvalSymlinks(dirName)

				mutex.Lock()
				pathMatches[dirName] = linkPath
				mutex.Unlock()

				if err != nil {
					fmt.Println(err)
				}
			} else {
				mutex.Lock()
				pathMatches[dirName] = ""
				mutex.Unlock()
			}
		}
	}
}

func findExecutable(nameToSearchFor string, noPath bool) {
	var searchPaths []string

	if !noPath {
		envPath := paths.GetAdditionalPaths()

		searchPaths = append(envPath, paths.SYSTEM_SEARCH_PATHS[0:8]...)
		searchPaths = slice.Unique(searchPaths)
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
	var args ProgramArgs

	arg.MustParse(&args)
	findExecutable(args.Name, args.NoPath)
	listMatches(args.Ugly)
}
