package version

import "runtime"

var (
	// Version is dynamically set by the toolchain or overridden by the Makefile.
	Version = "dev"
	// GoOs holds OS name.
	GoOs = runtime.GOOS
	// GoArch holds architecture name.
	GoArch = runtime.GOARCH

	// Windows, Linux, others os env var $KUBECONFIG splitter according to
	// https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/#append-home-kube-config-to-your-kubeconfig-environment-variable
	KubeConfigSplitter = map[string]string{
		Linux:   ":",
		Windows: ";",
		Others:  ":",
	}
)

const (
	Darwin  = "darwin"
	Linux   = "linux"
	Windows = "windows"
	FreeBSD = "freebsd"
	Others  = "others"
)
