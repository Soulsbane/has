package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
			//fmt.Println(path, info.Size())
			return nil
		})

	if err != nil {
		log.Println(err)
	}
}

func main() {
	args := os.Args[1:]
	//searchPath()
	if len(args) > 0 {
		nameOfProgram := args[0]

		for _, f := range searchPaths {
			searchPath(f, nameOfProgram)
		}
	}
}
