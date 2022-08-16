/*
Copyright © 2020 Guo Xudong

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
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sunny0826/kubecm/version"

	"github.com/cli/safeexec"
	"github.com/mgutz/ansi"
	"github.com/sunny0826/kubecm/pkg/update"

	"github.com/sunny0826/kubecm/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure" // required for Azure
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"   // required for GKE
)

func main() {
	kubecmVersion := version.Version

	updateMessageChan := make(chan *update.ReleaseInfo)
	go func() {
		rel, _ := update.CheckForUpdate("sunny0826/kubecm", kubecmVersion)
		updateMessageChan <- rel
	}()
	baseCommand := cmd.NewBaseCommand()
	if err := baseCommand.CobraCmd().Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	newRelease := <-updateMessageChan
	if newRelease != nil {
		fmt.Printf("\n\n%s %s → %s\n",
			ansi.Color("A new release of kubecm is available:", "yellow"),
			ansi.Color(strings.TrimPrefix(kubecmVersion, "v"), "cyan"),
			ansi.Color(strings.TrimPrefix(newRelease.Version, "v"), "green"))
		if isUnderHomebrew() {
			fmt.Printf("To upgrade, run: %s\n", "brew update && brew upgrade kubecm")
		}
		fmt.Printf("%s\n\n",
			ansi.Color(newRelease.URL, "yellow"))
	}
}

// Check whether the gh binary was found under the Homebrew prefix
func isUnderHomebrew() bool {
	brewExe, err := safeexec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return false
	}
	kubecmBinary, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(kubecmBinary, brewBinPrefix)
}
