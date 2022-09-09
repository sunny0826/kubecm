<p align="center">
    <img src="docs/static/kubecm.png" title="KubeCM" alt="Kubecm" height="200" />
</p>

![Go version](https://img.shields.io/github/go-mod/go-version/sunny0826/kubecm)
![Go](https://github.com/sunny0826/kubecm/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)
[![codecov](https://codecov.io/gh/sunny0826/kubecm/branch/master/graph/badge.svg?token=KGTLBQ8HYZ)](https://codecov.io/gh/sunny0826/kubecm)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/6065/badge)](https://bestpractices.coreinfrastructure.org/projects/6065)

```text
                                                 
        Manage your kubeconfig more easily.        
                                                   

██   ██ ██    ██ ██████  ███████  ██████ ███    ███ 
██  ██  ██    ██ ██   ██ ██      ██      ████  ████ 
█████   ██    ██ ██████  █████   ██      ██ ████ ██ 
██  ██  ██    ██ ██   ██ ██      ██      ██  ██  ██ 
██   ██  ██████  ██████  ███████  ██████ ██      ██

 Tips  Find more information at: https://kubecm.cloud

Usage:
  kubecm [command]

Available Commands:
  add         Add KubeConfig to $HOME/.kube/config
  alias       Generate alias for all contexts
  clear       Clear lapsed context, cluster and user
  cloud       Manage kubeconfig from cloud
  completion  Generate completion script
  create      Create new KubeConfig(experiment)
  delete      Delete the specified context from the kubeconfig
  help        Help about any command
  list        List KubeConfig
  merge       Merge multiple kubeconfig files into one
  namespace   Switch or change namespace interactively
  rename      Rename the contexts of kubeconfig
  switch      Switch Kube Context interactively
  version     Print version info

Flags:
      --config string   path of kubeconfig (default "$HOME/.kube/config")
  -h, --help            help for kubecm
      --ui-size int     number of list items to show in menu at once (default 4)

Use "kubecm [command] --help" for more information about a command.
```

## Documentation

For full documentation, please visit the KubeCM website: [https://kubecm.cloud](https://kubecm.cloud)

## Demo

[![asciicast](https://asciinema.org/a/389595.svg)](https://asciinema.org/a/389595)

## Install
Using [Krew](https://krew.sigs.k8s.io/):

```bash
kubectl krew install kc
```

Using Homebrew:

```bash
brew install kubecm
```

Source binary:

[Download the binary](https://github.com/sunny0826/kubecm/releases)

## Contribute

Feel free to open issues and pull requests. Any feedback is highly appreciated!

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=sunny0826/kubecm&type=Date)](https://star-history.com/#sunny0826/kubecm)


## Thanks

- [JetBrains IDEs](https://www.jetbrains.com/?from=kubecm)

<p align="center">
  <a href="https://www.jetbrains.com/?from=kubecm" title="前往官网了解JetBrains出品的IDEs">
    <img src="docs/static/jetbrains.svg" width="128" alt="JetBrains logo">
  </a>
</p>
