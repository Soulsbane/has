package main

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

	//"/usr/local/bin",
	//"/usr/local/sbin",
	//"/usr/share",

	//"/bin",
	//"/sbin",
	//	"/opt/bin",
}

func isValidPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func searchPath(path string, name string) {
	if isValidPath(path) {
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {

				if err != nil {
					return err
				}

				if filepath.Base(path) == name {
					dir := color.BlueString(filepath.Dir(path) + "/")
					base := color.GreenString(filepath.Base(path))
					fmt.Printf("%s%s", dir, base)
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
