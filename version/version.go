package version

import "runtime"

var (
	// Version is dynamically set by the toolchain or overridden by the Makefile.
	Version = "dev"
	// GitRevision is the commit of repo
	GitRevision = "UNKNOWN"
	// BuildDate is the build date of kubecm binary.
	BuildDate = ""
	// GoOs holds OS name.
	GoOs = runtime.GOOS
	// GoArch holds architecture name.
	GoArch = runtime.GOARCH
)
