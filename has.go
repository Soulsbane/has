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

func getEnvVarPaths() []string {
	path, variableExists := os.LookupEnv("PATH")
	var paths []string

	if variableExists {
		paths = strings.Split(path, ":")
	}

	return paths
}

// Taken from https://www.reddit.com/r/golang/comments/5ia523/idiomatic_way_to_remove_duplicates_in_a_slice/
func removeDuplicateDirs(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0

	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}

	return s[:j]
}

func searchDir(dirName string, nameToSearchFor string, de *godirwalk.Dirent) {
	if isValidPath(dirName) {
		if de.IsSymlink() && filepath.Base(dirName) == nameToSearchFor {
			linkPath, err := filepath.EvalSymlinks(dirName)

			fmt.Printf("Link: %s => %s\n", color.YellowString(dirName), color.BlueString(linkPath))
			dirName = linkPath

			if err != nil {
				fmt.Println(err)
			}
		}

		if filepath.Base(dirName) == nameToSearchFor && !de.IsDir() {
			dir := color.BlueString(filepath.Dir(dirName) + "/")
			base := color.GreenString(filepath.Base(dirName))
			fmt.Printf("%s%s\n", dir, base)
		}
	}
}

func searchDirs(nameToSearchFor string, noPathEnvVar bool) {
	if noPathEnvVar == false {
		searchPaths = append(searchPaths, getEnvVarPaths()...)
	}

	searchPaths = removeDuplicateDirs(searchPaths)

	for _, dirToSearch := range searchPaths {
		err := godirwalk.Walk(dirToSearch, &godirwalk.Options{
			Callback: func(walkDir string, de *godirwalk.Dirent) error {
				searchDir(walkDir, nameToSearchFor, de)
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
		FileName     string `arg:"positional, required"`
		NoPathEnvVar bool   `arg:"-n, --no-path" default:"false" help:"Include directories in user's $PATH."`
	}

	arg.MustParse(&args)
	searchDirs(args.FileName, args.NoPathEnvVar)
}
