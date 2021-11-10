module github.com/sunny0826/kubecm

go 1.15

require (
	github.com/alibabacloud-go/cs-20151215/v2 v2.4.5
	github.com/alibabacloud-go/darabonba-env v1.0.0
	github.com/alibabacloud-go/darabonba-openapi v0.1.7
	github.com/alibabacloud-go/tea v1.1.17
	github.com/alibabacloud-go/tea-utils v1.4.1 // indirect
	github.com/bndr/gotabulate v1.1.3-0.20170315142410-bc555436bfd5
	github.com/daviddengcn/go-colortext v1.0.0
	github.com/imdario/mergo v0.3.7
	github.com/manifoldco/promptui v0.3.2
	github.com/pterm/pterm v0.12.8
	github.com/rsteube/cobra-zsh-gen v1.1.0
	github.com/spf13/cobra v1.0.0
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v0.19.3
	k8s.io/utils v0.0.0-20201015054608-420da100c033 // indirect
)

replace github.com/manifoldco/promptui => github.com/terryding77/promptui v0.3.3
