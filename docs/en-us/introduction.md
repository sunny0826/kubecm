
List, switch, add, delete and more interactive operations to manage kubeconfig. 
It also supports kubeconfig management from cloud.

## Quick start

![demo](../static/kubecm-home.gif)

### 💫 Highlights

- **Context Management**: Switch between Kubernetes **clusters** and **namespaces** in a single command.
- **Merge-Kubeconfig**: Merge multiple kubeconfig files into one.
- **Interactive Mode**: Interactively select the context you want to switch to.
- **Multi-Platform**: Support Linux, macOS, and Windows.
- **Auto-Completion**: Support auto-completion for Bash, Zsh, and Fish.

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
![ns](../../static/ns-lates.gif)


### Interactive operation

<script src="https://asciinema.org/a/vL1vJZA1KeeFka9C0Wx3SSWto.js" id="asciicast-vL1vJZA1KeeFka9C0Wx3SSWto" async="true"></script>

more commands, please see [CLI References](./cli/kubecm_add.md)