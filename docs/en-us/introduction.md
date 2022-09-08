
List, switch, add, delete and more interactive operations to manage kubeconfig. 
It also supports kubeconfig management from cloud.

## Quick start

### Install

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

### Add kubeconfig

```bash
# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -
```

### List kubeconfig

```bash
# List all the contexts in your KubeConfig file
kubecm list
```

### Switch kubeconfig

```bash
# Switch Kube Context interactively
kubecm switch
# Quick switch Kube Context
kubecm switch dev
```

### Switch namespace

```bash
# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system
```
![ns](../../static/ns.gif)

### Interactive operation

<script id="asciicast-389595" src="https://asciinema.org/a/389595.js" async></script>

more commands, please see [CLI References](./cli/kubecm_add.md)