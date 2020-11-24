## kubecm completion

生成 bash/zsh 自动补全脚本

### 简介

输出自动补全脚本（bash 或 zsh）

```
kubecm completion [flags]
```

### 示例

```

# bash
kubecm completion bash > ~/.kube/kubecm.bash.inc
printf "
# kubecm shell completion
source '$HOME/.kube/kubecm.bash.inc'
" >> $HOME/.bash_profile
source $HOME/.bash_profile

# add to $HOME/.zshrc
source <(kubecm completion zsh)
# or
kubecm completion zsh > "${fpath[1]}/_kubecm"

```

### 选项

```
  -h, --help   help for completion
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
