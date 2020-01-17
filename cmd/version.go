/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
)

type VersionCommand struct {
	baseCommand
}

var (
	kubecmVersion = "unknown"
	goos          = runtime.GOOS
	goarch        = runtime.GOARCH
	gitCommit     = "$Format:%H$"          // sha1 from git, output of $(git rev-parse HEAD)
	buildDate     = "1970-01-01T00:00:00Z" // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

// version returns the version of kustomize.
type version struct {
	// KustomizeVersion is a kustomize binary version.
	kubecmVersion string `json:"kubecmVersion"`
	// GitCommit is a git commit
	GitCommit string `json:"gitCommit"`
	// BuildDate is a build date of the binary.
	BuildDate string `json:"buildDate"`
	// GoOs holds OS name.
	GoOs string `json:"goOs"`
	// GoArch holds architecture name.
	GoArch string `json:"goArch"`
}

func (vc *VersionCommand) Init() {
	vc.command = &cobra.Command{
		Use:   "version",
		Short: "Print version info",
		Long:  "Print version info",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version: %s\n", getVersion().kubecmVersion)
			cmd.Printf("GitCommit: %s\n", getVersion().GitCommit)
			cmd.Printf("BuildDate: %s\n", getVersion().BuildDate)
			cmd.Printf("GoOs: %s\n", getVersion().GoOs)
			cmd.Printf("GoArch: %s\n", getVersion().GoArch)
		},
	}
}

// getVersion returns version.
func getVersion() version {
	return version{
		kubecmVersion,
		gitCommit,
		buildDate,
		goos,
		goarch,
	}
}
