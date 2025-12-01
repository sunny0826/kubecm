## kubecm

KubeConfig Manager.

### Synopsis


[104m[104m                                                   [0m[0m
[104m[104m[97m[97m        Manage your kubeconfig more easily.        [0m[104m[0m[104m[0m[0m
[104m[104m                                                   [0m[0m
[92m[0m
[92m██   ██ ██    ██ ██████  ███████  ██████ ███    ███ [0m
[92m██  ██  ██    ██ ██   ██ ██      ██      ████  ████ [0m
[92m█████   ██    ██ ██████  █████   ██      ██ ████ ██ [0m
[92m██  ██  ██    ██ ██   ██ ██      ██      ██  ██  ██ [0m
[92m██   ██  ██████  ██████  ███████  ██████ ██      ██[0m
[92m[0m
[44;97m[44;97m Tips [0m[0m [96m[96mFind more information at: ]8;;https://kubecm.cloud[3;32mkubecm.cloud]8;;[0m[96m[0m[0m


### Options

```
      --config string   path of kubeconfig (default "$HOME/.kube/config")
      --create          Create a new kubeconfig file if not exists
  -h, --help            help for kubecm
  -m, --mac-notify      enable to display Mac notification banner
  -s, --silence-table   enable/disable output of context table on successful config update
  -u, --ui-size int     number of list items to show in menu at once (default 10)
```

### SEE ALSO

* [kubecm add](kubecm_add.md)	 - Add KubeConfig to $HOME/.kube/config
* [kubecm alias](kubecm_alias.md)	 - Generate alias for all contexts
* [kubecm clear](kubecm_clear.md)	 - Clear lapsed context, cluster and user
* [kubecm cloud](kubecm_cloud.md)	 - Manage kubeconfig from cloud
* [kubecm create](kubecm_create.md)	 - Create new KubeConfig(experiment)
* [kubecm delete](kubecm_delete.md)	 - Delete the specified context from the kubeconfig
* [kubecm docs](kubecm_docs.md)	 - Open document website
* [kubecm export](kubecm_export.md)	 - Export the specified context from the kubeconfig
* [kubecm list](kubecm_list.md)	 - List KubeConfig
* [kubecm merge](kubecm_merge.md)	 - Merge multiple kubeconfig files into one
* [kubecm rename](kubecm_rename.md)	 - Rename the contexts of kubeconfig
* [kubecm switch](kubecm_switch.md)	 - Switch Kube Context interactively
* [kubecm version](kubecm_version.md)	 - Print version info

