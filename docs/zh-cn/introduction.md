
List, switch, add, delete and more interactive operations to manage kubeconfig.
It also supports kubeconfig management from cloud.

## 快速开始

### 安装

使用 [Krew](https://krew.sigs.k8s.io/):

```bash
kubectl krew install kc
```

使用 Homebrew:

```bash
brew install kubecm
```

Source binary:

[下载可执行文件](https://github.com/sunny0826/kubecm/releases)

### 添加 kubeconfig

```bash
# Merge test.yaml with $HOME/.kube/config
kubecm add -f test.yaml 
# Add kubeconfig from stdin
cat /etc/kubernetes/admin.conf |  kubecm add -f -
```

### 列举 kubeconfig

```bash
# List all the contexts in your KubeConfig file
kubecm list
```

### 切换 kubeconfig

```bash
# Switch Kube Context interactively
kubecm switch
# Quick switch Kube Context
kubecm switch dev
```

### 切换 namespace

```bash
# Switch Namespace interactively
kubecm namespace
# or
kubecm ns
# change to namespace of kube-system
kubecm ns kube-system
```
![ns](../../static/ns.gif)

## 交互式操作

<script id="asciicast-389595" src="https://asciinema.org/a/389595.js" async></script>

更多的信息, 请看 [CLI 参考](https://kubecm.cloud/zh-cn/cli/kubecm_add)