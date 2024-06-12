Multiple installation paths are available.

<!-- tabs:start -->

#### ** Krew **

[![homebrew downloads](https://img.shields.io/homebrew/installs/dy/kubecm?style=for-the-badge&logo=homebrew&label=homebrew)](https://formulae.brew.sh/formula/kubecm)
[![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/sunny0826/kubecm/total?style=for-the-badge&logo=github&label=github%20download)](https://github.com/sunny0826/kubecm/releases)

Using [Krew](https://krew.sigs.k8s.io/):

```bash
kubectl krew install kc
```

#### ** Homebrew **

```bash
brew install kubecm
```

#### ** Download the binary **

[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)

```bash
VERSION=v0.27.1 # replace with the version you want, note the "v" prefix!
# linux x86_64
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/${VERSION}/kubecm_${VERSION}_Linux_x86_64.tar.gz
# macos
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/${VERSION}/kubecm_${VERSION}_Darwin_x86_64.tar.gz
# windows
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/${VERSION}/kubecm_${VERSION}_Windows_x86_64.tar.gz

# linux & macos
tar -zxvf kubecm.tar.gz kubecm
cd kubecm
sudo mv kubecm /usr/local/bin/

# windows
# Unzip kubecm.tar.gz
# Add the binary in to your $PATH
```

<!-- tabs:end -->