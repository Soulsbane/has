package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/alexflint/go-arg"
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

func searchPath(path string, name string) {
	fmt.Println("Searching path: ", path)
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if filepath.Base(path) == name {
				fmt.Println("FOUND: ", path)
			}

			return nil
		})

	if err != nil {
		log.Println(err)
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
