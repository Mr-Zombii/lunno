package version

import "runtime"

var (
	Version   = "unknown"
	GitCommit = "unknown"
	GoVersion = runtime.Version()
	OS        = runtime.GOOS
	Arch      = runtime.GOARCH
)

func Full() string {
	return Version + " (" + GitCommit + ") [" + OS + "/" + Arch + " | " + GoVersion + "]"
}
