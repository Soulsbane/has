package paths

import (
	"os"
	"strings"
)

func GetAdditionalPaths() []string {
	path, variableExists := os.LookupEnv("PATH")
	var paths []string

	if variableExists {
		paths = strings.Split(path, ":")
	}

	return paths
}
