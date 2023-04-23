## kubecm switch

交互式切换 Kube Context

### 简介


切换 Kube Context，可以直接切换或交互式切换。

**可以通过键入 `\` 来进行搜索**

```
kubecm switch [flags]
```

![switch](../../static/switch.gif)

### 示例

```

# Switch Kube Context interactively
kubecm switch
# Quick switch Kube Context
kubecm switch dev

```

### 选项

```
  -h, --help   help for switch
```

### 全局选项

```
      --config string   path of kubeconfig (default "/Users/guoxudong/.kube/config")
  -m, --mac-notify      enable to display Mac notification banner
      --ui-size int     number of list items to show in menu at once (default 4)
```
