## kubecm alias

Generate alias for all contexts

### Synopsis

Generate alias for all contexts

```
kubecm alias
```

### Examples

```

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
      --ui-size int     number of list items to show in menu at once (default 4)
```

### SEE ALSO

* [kubecm](kubecm.md)	 - KubeConfig Manager.

