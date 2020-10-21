[English](../README.md) | 简体中文

# KubeCM

[![Build Status](https://travis-ci.org/sunny0826/kubecm.svg?branch=master)](https://travis-ci.org/sunny0826/kubecm)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)

```text
KubeConfig Manager
 _          _
| | ___   _| |__   ___  ___ _ __ ___
| |/ / | | | '_ \ / _ \/ __| '_ \ _ \
|   <| |_| | |_) |  __/ (__| | | | | |
|_|\_\\__,_|_.__/ \___|\___|_| |_| |_|

Find more information at: https://github.com/sunny0826/kubecm

Usage:
  kubecm [flags]
  kubecm [command]

Available Commands:
    add         Merge configuration file with $HOME/.kube/config
    alias       Generate alias for all contexts
    completion  Generates bash/zsh completion scripts
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

## 快速开始

### 安装

#### Homebrew

```bash
brew install sunny0826/tap/kubecm
```

#### 下载可执行文件

[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)

```bash
# linux x86_64
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/v${VERSION}/kubecm_${VERSION}_Linux_x86_64.tar.gz
# macos
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/v${VERSION}/kubecm_${VERSION}_Darwin_x86_64.tar.gz
# windows
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/v${VERSION}/kubecm_${VERSION}_Windows_x86_64.tar.gz

# linux & macos
tar -zxvf kubecm.tar.gz kubecm
cd kubecm
sudo mv kubecm /usr/local/bin/

# windows
# Unzip kubecm.tar.gz
# Add the binary in to your $PATH
```

### 自动补全

#### bash

```bash
# bash
kubecm completion bash > ~/.kube/kubecm.bash.inc
printf "
# kubecm shell completion
source '$HOME/.kube/kubecm.bash.inc'
" >> $HOME/.bash_profile
source $HOME/.bash_profile
```

#### zsh

```bash
# add to $HOME/.zshrc 
source <(kubecm completion zsh)
# or
kubecm completion zsh > "${fpath[1]}/_kubecm"
```

### 交互式操作

![Interactive](Interaction.gif)

### 列举所有 Context

```shell script
# List all the contexts in your kubeconfig file
kubecm ls
# Aliases
kubecm l
```

### 生成 alias

为所有 contexts 生成 alias，可自动将其注入到 `.zshrc` 或 `.bash_profile` 中。

```shell script
$ kubecm alias
# dev 
alias k-dev='kubectl --context dev'
# test
alias k-test='kubectl --context test'
# prod
alias k-prod='kubectl --context prod'
$ kubecm alias -o zsh
# add alias context to ~/.zshrc
$ kubecm alias -o bash
# add alias context to ~/.bash_profile
```


### 将 kubeconfig 加入到 `$HOME/.kube/config`

```shell script
# Merge example.yaml with $HOME/.kube/config.yaml
kubecm add -f example.yaml 

# Merge example.yaml and name contexts test with $HOME/.kube/config.yaml
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c
```

### 合并 kubeconfig

```shell script
# Merge kubeconfig in the directory
kubecm merge -f dir

# Merge kubeconfig in the directory and overwrite the original kubeconfig file
kubecm merge -f dir -c
```

###  交互式切换 Kube Context

```shell script
# Switch Kube Context interactively
kubecm switch
```
![switch](switch.gif)

### 删除 context

```shell script
# Delete the context interactively
kubecm delete
# Delete the context
kubecm delete my-context
```
### 重命名 context

```shell script
# Renamed the context interactively
kubecm rename
```
 
### 切换当前集群 namespace

You can switch namespace or switch namespace interactively

```shell script
# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system
``` 
![ns](ns.gif)

## 视频

[![](https://tva3.sinaimg.cn/large/ad5fbf65gy1gij1pl0pn5j218o0p81kx.jpg)](https://www.bilibili.com/video/av88259938/)

## 贡献

欢迎在 issue 中与作者交流或提交 pull requests，您的反馈对本项目的健康发展是十分重要的。

## 鸣谢

- [JetBrains IDEs](https://www.jetbrains.com/?from=kubecm)

<p align="center">
  <a href="https://www.jetbrains.com/?from=kubecm" title="前往官网了解JetBrains出品的IDEs">
    <img src="dosc/jetbrains.svg" width="128" alt="JetBrains logo">
  </a>
</p>
