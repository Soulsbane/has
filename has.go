package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
	"github.com/karrick/godirwalk"
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

func searchDir(dirName string, nameToSearchFor string, de *godirwalk.Dirent) {
	if isValidPath(dirName) {
		if de.IsSymlink() && filepath.Base(dirName) == nameToSearchFor {
			linkPath, err := filepath.EvalSymlinks(dirName)
			dirName = linkPath
			pathMatches[dirName] = linkPath

			if err != nil {
				fmt.Println(err)
			}
		}

		if filepath.Base(dirName) == nameToSearchFor && !de.IsDir() {
			pathMatches[dirName] = ""
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

func searchDirs(nameToSearchFor string, noPath bool) {
	if !noPath {
		path, _ := exec.LookPath(nameToSearchFor)
		lookPath(path)
	}

	for _, dirToSearch := range searchPaths {
		err := godirwalk.Walk(dirToSearch, &godirwalk.Options{
			Callback: func(walkDir string, de *godirwalk.Dirent) error {
				searchDir(walkDir, nameToSearchFor, de)
				return nil
			},
			Unsorted: true,
		})

		if err != nil {
			// FIXME: Permission errors from distros using wrong permissions.
		}
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
	searchDirs(args.FileName, args.NoPath)
	listMatches(args.Ugly)
}
