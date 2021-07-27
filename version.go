package ec2fzf

import "fmt"

// Values for these are injected by the build.
var (
	BuildTime = "unset" // BuildTime is a time label of the moment when the binary was built
	Commit    = "unset" // Commit is a last commit hash at the moment when the binary was built
	Release   = "unset" // Release is a semantic Version of current build
	Arch      = "unset" // Arch type
)

func showVersion() {
	v := fmt.Sprintf("hello ec2-fzf\nRelease: %s\n"+
		"Commit: %s\n"+
		"Build Time: %s\n"+
		"Arch: %s",
		Release, Commit, BuildTime, Arch)

	fmt.Println(v)
}

// https://github.com/dapr/dapr/blob/master/Makefile
