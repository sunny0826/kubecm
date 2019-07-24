# KubeCM
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

# Merge example.yaml with ./kube/config
kubecm add -f example.yaml 

# Merge example.yaml and name contexts test with ./kube/config
kubecm add -f example.yaml -n test

# Overwrite the original kubeconfig file
kubecm add -f example.yaml -c


Flags:
  -c, --cover         Overwrite the original kubeconfig file
  -f, --file string   Path to merge kubeconfig files
  -h, --help          help for add
  -n, --name string   The name of contexts. if this field is null,it will be named with file name.

```