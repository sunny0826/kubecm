package version

import "runtime"

var (
	// Version is dynamically set by the toolchain or overridden by the Makefile.
	Version = "dev"
	// GoOs holds OS name.
	GoOs = runtime.GOOS
	// GoArch holds architecture name.
	GoArch = runtime.GOARCH
)
