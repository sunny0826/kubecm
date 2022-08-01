# KubeCM

![Go](https://github.com/sunny0826/kubecm/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)
[![codecov](https://codecov.io/gh/sunny0826/kubecm/branch/master/graph/badge.svg?token=KGTLBQ8HYZ)](https://codecov.io/gh/sunny0826/kubecm)

```text
Manage your kubeconfig more easily.

██   ██ ██    ██ ██████  ███████  ██████ ███    ███
██  ██  ██    ██ ██   ██ ██      ██      ████  ████ 
█████   ██    ██ ██████  █████   ██      ██ ████ ██ 
██  ██  ██    ██ ██   ██ ██      ██      ██  ██  ██ 
██   ██  ██████  ██████  ███████  ██████ ██      ██

Find more information at: https://github.com/sunny0826/kubecm

Usage:
  kubecm [flags]
  kubecm [command]

Available Commands:
    add         Merge configuration file with $HOME/.kube/config
    alias       Generate alias for all contexts
    clear       Clear lapsed context, cluster and user
    completion  Generates bash/zsh completion scripts
    create      Create new KubeConfig(experiment)
    delete      Delete the specified context from the kubeconfig
    help        Help about any command
    ls          List kubeconfig
    merge       Merge the kubeconfig files in the specified directory
    namespace   Switch or change namespace interactively
    rename      Rename the contexts of kubeconfig
    switch      Switch Kube Context interactively
    version     Print version info


Flags:
      --config string   path of kubeconfig (default "$HOME/.kube/config")
  -h, --help   help for kubecm

Use "kubecm [command] --help" for more information about a command.
```

# 视频

<script id="asciicast-389595" src="https://asciinema.org/a/389595.js" async></script>

# 鸣谢

- [JetBrains IDEs](https://www.jetbrains.com/?from=kubecm)

<p align="center">
  <a href="https://www.jetbrains.com/?from=kubecm" title="前往官网了解JetBrains出品的IDEs">
    <img src="../static/jetbrains.svg" width="128" alt="JetBrains logo">
  </a>
</p>
