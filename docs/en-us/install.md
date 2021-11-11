# Install

Multiple installation paths are available.

<!-- tabs:start -->

#### ** Homebrew **

```bash
brew install kubecm
```

#### ** Download the binary **

[![GitHub release](https://img.shields.io/github/release/sunny0826/kubecm)](https://github.com/sunny0826/kubecm/releases)

```bash
# linux x86_64
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/v${VERSION}/kubecm_${VERSION}_Linux_x86_64.tar.gz
# macos
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/v${VERSION}/kubecm_${VERSION}_Darwin_x86_64.tar.gz
# windows
curl -Lo kubecm.tar.gz https://github.com/sunny0826/kubecm/releases/download/v${VERSION}/kubecm_${VERSION}_Windows_x86_64.tar.gz

# linux & macos
tar -zxvf kubecm.tar.gz kubecm
cd kubecm
sudo mv kubecm /usr/local/bin/

# windows
# Unzip kubecm.tar.gz
# Add the binary in to your $PATH
```

<!-- tabs:end -->