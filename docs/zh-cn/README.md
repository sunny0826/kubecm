<p align="center">
    <img src="../../static/kubecm.png" title="KubeCM" alt="Kubecm" height="200" />
</p>

![Go version](https://img.shields.io/github/go-mod/go-version/sunny0826/kubecm)
![Go](https://github.com/sunny0826/kubecm/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)
[![codecov](https://codecov.io/gh/sunny0826/kubecm/branch/master/graph/badge.svg?token=KGTLBQ8HYZ)](https://codecov.io/gh/sunny0826/kubecm)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/6065/badge)](https://bestpractices.coreinfrastructure.org/projects/6065)

List, switch, add, delete and more interactive operations to manage kubeconfig. It also supports kubeconfig management from cloud.

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

## 演示

<script id="asciicast-389595" src="https://asciinema.org/a/389595.js" async></script>

## 安装
使用 [Krew](https://krew.sigs.k8s.io/):

```bash
kubectl krew install kc
```

使用 Homebrew:

```bash
brew install kubecm
```

下载可执行文件:

[Download the binary](https://github.com/sunny0826/kubecm/releases)

## 鸣谢

- [JetBrains IDEs](https://www.jetbrains.com/?from=kubecm)

<p align="center">
  <a href="https://www.jetbrains.com/?from=kubecm" title="前往官网了解JetBrains出品的IDEs">
    <img src="../static/jetbrains.svg" width="128" alt="JetBrains logo">
  </a>
</p>
