package version

import (
	"fmt"
	"runtime"
)

const binaryName string = "termin8"

var (
	version   = "notSet"
	buildTime = "1970-01-01T00:00:00Z"
	gitCommit = "notSet"
)

func PrintVersion() string {
	version := fmt.Sprintf("%s v%s", binaryName, version)
	return version
}

func GetBinaryName() string {
	return binaryName
}
func GetGitCommit() string {
	return gitCommit
}

func GetBuildTime() string {
	return buildTime
}

func GetGoVersion() string {
	return runtime.Version()
}

func GetGoPlatform() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

func GetGoCompiler() string {
	return runtime.Compiler
}
