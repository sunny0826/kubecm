# KubeCM

[![Build Status](https://travis-ci.org/sunny0826/kubecm.svg?branch=master)](https://travis-ci.org/sunny0826/kubecm)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)

Merge multiple kubeconfig

```bash
Merge multiple kubeconfig
 _          _
| | ___   _| |__   ___  ___ _ __ ___
| |/ / | | | '_ \ / _ \/ __| '_ \ _ \
|   <| |_| | |_) |  __/ (__| | | | | |
|_|\_\\__,_|_.__/ \___|\___|_| |_| |_|

Find more information at: https://github.com/sunny0826/kubecm

Usage:
  kubecm [command]

Available Commands:
  add         Merge configuration file with ./kube/config
  delete      Delete the specified context from the kubeconfig
  get         Displays one or many contexts from the kubeconfig file.
  help        Help about any command
  use         Sets the current-context in a kubeconfig file
  version     Prints the kubecm version

Flags:
  -h, --help   help for kubecm

Use "kubecm [command] --help" for more information about a command.
```

## Quick Start

### Install kubecm

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

tar -zxvf kubecm.tar.gz kubecm
sudo mv kubecm /usr/local/bin/
```

### Add configuration to `./kube/config`

```bash
# Merge example.yaml with ./kube/config.yaml
kubecm add -f example.yaml 

# Merge example.yaml and name contexts test with ./kube/config.yaml
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c
```

### Displays contexts

```bash
# List all the contexts in your kubeconfig file
kubecm get

# Describe one context in your kubeconfig file.
kubecm get my-context

# example output
$ kubecm get
+------------+-----------------------+---------------------------+------------------------+
|   CURRENT  |          NAME         |          CLUSTER          |          USER          |
+============+=======================+===========================+========================+
|      *     |      al_devops        |    al_devops-0-cluster    |    al_devops-0-user    |
+------------+-----------------------+---------------------------+------------------------+
|            |       al_prod         |     al_prod-0-cluster     |     al_prod-0-user     |
+------------+-----------------------+---------------------------+------------------------+
|            |       al_test         |     al_test-0-cluster     |     al_test-0-user     |
+------------+-----------------------+---------------------------+------------------------+

```

### Delete context

```bash
# Delete the context
kubecm delete my-context
```

### Set context

```bash
# Use the context for the test cluster
kubecm use test
```

## Contribute

Feel free to open issues and pull requests. Any feedback is highly appreciated!