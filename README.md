<p align="center">
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-3-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
    <img src="docs/static/kubecm.png" title="KubeCM" alt="Kubecm" height="200" />
</p>

![Go version](https://img.shields.io/github/go-mod/go-version/sunny0826/kubecm)
![Go](https://github.com/sunny0826/kubecm/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sunny0826/kubecm)](https://goreportcard.com/report/github.com/sunny0826/kubecm)
![GitHub](https://img.shields.io/github/license/sunny0826/kubecm.svg)
[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)
[![codecov](https://codecov.io/gh/sunny0826/kubecm/branch/master/graph/badge.svg?token=KGTLBQ8HYZ)](https://codecov.io/gh/sunny0826/kubecm)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/6065/badge)](https://bestpractices.coreinfrastructure.org/projects/6065)

```text
                                                 
        Manage your kubeconfig more easily.        
                                                   

██   ██ ██    ██ ██████  ███████  ██████ ███    ███ 
██  ██  ██    ██ ██   ██ ██      ██      ████  ████ 
█████   ██    ██ ██████  █████   ██      ██ ████ ██ 
██  ██  ██    ██ ██   ██ ██      ██      ██  ██  ██ 
██   ██  ██████  ██████  ███████  ██████ ██      ██

 Tips  Find more information at: https://kubecm.cloud

Usage:
  kubecm [command]

Available Commands:
  add         Add KubeConfig to $HOME/.kube/config
  alias       Generate alias for all contexts
  clear       Clear lapsed context, cluster and user
  cloud       Manage kubeconfig from cloud
  completion  Generate completion script
  create      Create new KubeConfig(experiment)
  delete      Delete the specified context from the kubeconfig
  help        Help about any command
  list        List KubeConfig
  merge       Merge multiple kubeconfig files into one
  namespace   Switch or change namespace interactively
  rename      Rename the contexts of kubeconfig
  switch      Switch Kube Context interactively
  version     Print version info

Flags:
      --config string   path of kubeconfig (default "$HOME/.kube/config")
  -h, --help            help for kubecm
      --ui-size int     number of list items to show in menu at once (default 4)

Use "kubecm [command] --help" for more information about a command.
```

## Documentation

For full documentation, please visit the KubeCM website: [https://kubecm.cloud](https://kubecm.cloud)

## Demo

[![asciicast](https://asciinema.org/a/389595.svg)](https://asciinema.org/a/389595)

## Install
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

## Contribute

Feel free to open issues and pull requests. Any feedback is highly appreciated! I'd be grateful if you can give me some feedback.

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://guoxudong.io/"><img src="https://avatars.githubusercontent.com/u/24563928?v=4?s=100" width="100px;" alt="Xudong Guo"/><br /><sub><b>Xudong Guo</b></sub></a><br /><a href="https://github.com/sunny0826/kubecm/commits?author=sunny0826" title="Code">💻</a> <a href="https://github.com/sunny0826/kubecm/commits?author=sunny0826" title="Documentation">📖</a> <a href="#maintenance-sunny0826" title="Maintenance">🚧</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/astraw99"><img src="https://avatars.githubusercontent.com/u/25424734?v=4?s=100" width="100px;" alt="Cheng Wang"/><br /><sub><b>Cheng Wang</b></sub></a><br /><a href="https://github.com/sunny0826/kubecm/commits?author=astraw99" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://nolebase.ayaka.io/"><img src="https://avatars.githubusercontent.com/u/11081491?v=4?s=100" width="100px;" alt="Neko Ayaka"/><br /><sub><b>Neko Ayaka</b></sub></a><br /><a href="https://github.com/sunny0826/kubecm/commits?author=nekomeowww" title="Code">💻</a></td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td align="center" size="13px" colspan="7">
        <img src="https://raw.githubusercontent.com/all-contributors/all-contributors-cli/1b8533af435da9854653492b1327a23a4dbd0a10/assets/logo-small.svg">
          <a href="https://all-contributors.js.org/docs/en/bot/usage">Add your contributions</a>
        </img>
      </td>
    </tr>
  </tfoot>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=sunny0826/kubecm&type=Date)](https://star-history.com/#sunny0826/kubecm)

## Thanks

- [JetBrains IDEs](https://www.jetbrains.com/?from=kubecm)

<p align="center">
  <a href="https://www.jetbrains.com/?from=kubecm" title="前往官网了解JetBrains出品的IDEs">
    <img src="docs/static/jetbrains.svg" width="128" alt="JetBrains logo">
  </a>
</p>

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!