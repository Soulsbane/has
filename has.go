package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func isValidPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func isValidLinkPath(info os.FileInfo, path string) string {
	mode := info.Mode()
	link := mode & os.ModeSymlink

	if link != 0 {
		linkPath, err := filepath.EvalSymlinks(path)

		if err != nil {
			return ""
		}

		return linkPath
	}

	return ""
}

func getEnvVarPaths() []string {
	path, variableExists := os.LookupEnv("PATH")
	var paths []string

	if variableExists {
		paths = strings.Split(path, ":")
	}

	return paths
}

func searchPath(path string, name string) {
	if isValidPath(path) {
		err := godirwalk.Walk(path, &godirwalk.Options{
			Callback: func(currentPath string, de *godirwalk.Dirent) error {
				if isValidPath(currentPath) {
					if de.IsSymlink() {
						linkPath, err := filepath.EvalSymlinks(currentPath)
						currentPath = linkPath

						if err != nil {
							fmt.Println(err)
						}
					}

					if filepath.Base(currentPath) == name && !de.IsDir() {
						dir := color.BlueString(filepath.Dir(currentPath) + "/")
						base := color.GreenString(filepath.Base(currentPath))
						fmt.Printf("%s%s\n", dir, base)
					}
				}
				return nil
			},
			Unsorted: true,
		})

		if err != nil {
			// FIXME: Permission errors still linger so silence them for now.
		}
	}
}

func main() {
	var args struct {
		FileName      string `arg:"positional, required"`
		UsePathEnvVar bool   `arg:"-p, --path" default:"false" help:"Include directories in user's $PATH."`
	}

	arg.MustParse(&args)

	if args.UsePathEnvVar {
		searchPaths = append(searchPaths, getEnvVarPaths()...)
	}

	for _, f := range searchPaths {
		searchPath(f, args.FileName)
	}
}
