package main

// TODO: Perhaps a feature to use user's path also?

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
)

var searchPaths = [...]string{
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
		FileName string `arg:"positional, required"`
	}

	arg.MustParse(&args)

	for _, f := range searchPaths {
		searchPath(f, args.FileName)
	}
}
