# KubeCM

[![Build Status](https://travis-ci.org/sunny0826/kubecm.svg?branch=master)](https://travis-ci.org/sunny0826/kubecm)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)

Merge multiple kubeconfig

```bash
 _          _
| | ___   _| |__   ___  ___ _ __ ___
| |/ / | | | '_ \ / _ \/ __| '_ \ _ \
|   <| |_| | |_) |  __/ (__| | | | | |
|_|\_\\__,_|_.__/ \___|\___|_| |_| |_|

Usage:
  kubecm [command]

Available Commands:
  add         Merge configuration file with ./kube/config
  help        Help about any command

Flags:
  -h, --help   help for kubecm

Use "kubecm [command] --help" for more information about a command.
---------------------------------------------------------------------------
Merge configuration file with ./kube/config

Usage:
  kubecm add [flags]

Examples:

config.yaml
kubecm add -f example.yaml 

config.yaml
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c


Flags:
  -c, --cover         Overwrite the original kubeconfig file
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
  -n, --name string   The name of contexts. if this field is null,it will be named with file name.

```