## kubecm completion

Generates bash/zsh completion scripts

### Synopsis

Output shell completion code for the specified shell (bash or zsh).

```
kubecm completion [flags]
```

### Examples

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

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
