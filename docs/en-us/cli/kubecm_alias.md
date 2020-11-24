## kubecm alias

Generate alias for all contexts

### Synopsis

Generate alias for all contexts

```
kubecm alias
```

### Examples

```

# dev 
alias k-dev='kubectl --context dev'
# test
alias k-test='kubectl --context test'
# prod
alias k-prod='kubectl --context prod'
$ kubecm alias -o zsh
# add context to ~/.zshrc
$ kubecm alias -o bash
# add context to ~/.bash_profile

```

### Options

```
  -h, --help         help for alias
  -o, --out string   output to ~/.zshrc or ~/.bash_profile
```

### Options inherited from parent commands

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
```
