package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
)

var searchPaths = []string{
	"/usr/bin",
	//"/usr/lib",
	"/usr/sbin",

	"/usr/local/bin",
	"/usr/local/sbin",
	//"/usr/share", // Needs permission

	"/bin",
	"/sbin",
	"/opt/bin",
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
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {

				if err != nil {
					return err
				}

				isLinkPath := isValidLinkPath(info, path)

				if isLinkPath != "" {
					path = isLinkPath
				}

				if filepath.Base(path) == name && !info.IsDir() {
					dir := color.BlueString(filepath.Dir(path) + "/")
					base := color.GreenString(filepath.Base(path))
					fmt.Printf("%s%s\n", dir, base)
				}

				return nil
			})

		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	var args struct {
		FileName      string `arg:"positional, required"`
		UsePathEnvVar bool   `arg:"-u" default:"false" help:"Include directories in user's $PATH."`
	}

	arg.MustParse(&args)

	if args.UsePathEnvVar {
		searchPaths = append(searchPaths, getEnvVarPaths()...)
	}

	for _, f := range searchPaths {
		searchPath(f, args.FileName)
	}
}
