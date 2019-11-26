# KubeCM

[![Build Status](https://travis-ci.org/sunny0826/kubecm.svg?branch=master)](https://travis-ci.org/sunny0826/kubecm)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)

KubeConfig Manager

```bash
KubeConfig Manager
 _          _
| | ___   _| |__   ___  ___ _ __ ___
| |/ / | | | '_ \ / _ \/ __| '_ \ _ \
|   <| |_| | |_) |  __/ (__| | | | | |
|_|\_\\__,_|_.__/ \___|\___|_| |_| |_|

Find more information at: https://github.com/sunny0826/kubecm

Usage:
  kubecm [command]

Available Commands:
  add         Merge configuration file with $HOME/.kube/config
  delete      Delete the specified context from the kubeconfig
  get         Displays one or many contexts from the kubeconfig file.
  help        Help about any command
  merge       Merge the kubeconfig files in the specified directory.
  rename      Rename the contexts of kubeconfig
  use         Sets the current-context in a kubeconfig file
  version     Prints the kubecm version


Flags:
  -h, --help   help for kubecm

Use "kubecm [command] --help" for more information about a command.
```

## Quick Start

### Install

Homebrew:

```bash
brew install sunny0826/tap/kubecm
```

Download the binary:

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

### Add configuration to `$HOME/.kube/config`

```bash
# Merge example.yaml with $HOME/.kube/config.yaml
kubecm add -f example.yaml 

# Merge example.yaml and name contexts test with $HOME/.kube/config.yaml
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c
```

### Merge the kubeconfig

```bash
# Merge kubeconfig in the test directory
kubecm merge -f test 

# Merge kubeconfig in the test directory and overwrite the original kubeconfig file
kubecm merge -f test -c
```

### Displays contexts

```bash
# List all the contexts in your kubeconfig file
kubecm get

# Describe one context in your kubeconfig file.
kubecm get my-context

# example output
$ kubecm get
+------------+-----------------------+-----------------------+--------------------+--------------+
|   CURRENT  |          NAME         |        CLUSTER        |        USER        |   Namespace  |
+============+=======================+=======================+====================+==============+
|      *     |         test          |   cluster-28989kd95m  |   user-28989kd95m  |              |
+------------+-----------------------+-----------------------+--------------------+--------------+
|            |        test-1         |   cluster-7thmtkbk6m  |   user-7thmtkbk6m  |              |
+------------+-----------------------+-----------------------+--------------------+--------------+
|            |        test-2         |   cluster-4h9m74h8d6  |   user-4h9m74h8d6  |              |
+------------+-----------------------+-----------------------+--------------------+--------------+
```

### Delete context

```bash
# Delete the context
kubecm delete my-context
```

### Switch context

```bash
# Use the context for the test cluster
kubecm use test
```

### Rename context

```bash
# Renamed dev to test
kubecm rename -o dev -n test
# Renamed current-context name to dev
kubecm rename -n dev -c
```

## Contribute

Feel free to open issues and pull requests. Any feedback is highly appreciated!

## Thanks

- [JetBrains IDEs](https://www.jetbrains.com/?from=kubecm)

<p align="center">
  <a href="https://www.jetbrains.com/?from=kubecm" title="前往官网了解JetBrains出品的IDEs">
    <img src="dosc/jetbrains.svg" width="128" alt="JetBrains logo">
  </a>
</p>
