package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// CreateCommand clean command struct
type CreateCommand struct {
	BaseCommand
}

type CreateOptions struct {
	config      *clientcmdapi.Config
	clientSet   *kubernetes.Clientset
	role        string
	token       string
	contextName string
	userName    string
	namespace   string
}

// Init CreateCommand
func (ce *CreateCommand) Init() {
	ce.command = &cobra.Command{
		Use:   "create",
		Short: "Create new KubeConfig(experiment)",
		Long:  "Create new KubeConfig(experiment)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ce.runCreate(cmd, args)
		},
		Example: createExample(),
	}
	ce.command.DisableFlagsInUseLine = true
}

func (ce *CreateCommand) runCreate(cmd *cobra.Command, args []string) error {
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	userName := PromptUI("user name", "")
	co := CreateOptions{
		config:   config,
		userName: userName,
	}
	err = co.chooseContext()
	if err != nil {
		return err
	}
	err = co.chooseNamespace()
	if err != nil {
		return err
	}
	err = co.createServiceAccounts()
	if err != nil {
		return err
	}
	err = co.selectClusterRole()
	if err != nil {
		return err
	}
	err = co.createRoleBinding()
	if err != nil {
		return err
	}
	err = co.getToken()
	if err != nil {
		return err
	}
	newConfig := co.putOutKubeConfig()
	fileName := fmt.Sprintf("%s.kubeconfig", co.userName)
	err = clientcmd.WriteToFile(*newConfig, fileName)
	if err != nil {
		return err
	}
	printString(os.Stdout, "kubeconfig: "+fileName+" create success\n")
	return nil
}

func (co *CreateOptions) chooseContext() error {
	var kubeItems []Needle
	current := co.config.CurrentContext
	for key, obj := range co.config.Contexts {
		if key != current {
			kubeItems = append(kubeItems, Needle{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo})
		} else {
			kubeItems = append([]Needle{{Name: key, Cluster: obj.Cluster, User: obj.AuthInfo, Center: "(*)"}}, kubeItems...)
		}
	}
	num := SelectUI(kubeItems, "Select Kube Context")
	co.contextName = kubeItems[num].Name
	co.config.CurrentContext = co.contextName
	clientConfig := clientcmd.NewDefaultClientConfig(
		*co.config,
		&clientcmd.ConfigOverrides{},
	)
	c, _ := clientConfig.ClientConfig()
	clientSet, err := kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}
	co.clientSet = clientSet
	return nil
}

func (co *CreateOptions) chooseNamespace() error {
	var nss []Namespaces
	ctx := context.TODO()
	namespaceList, err := co.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, specItem := range namespaceList.Items {
		nss = append(nss, Namespaces{Name: specItem.Name, Default: false})
	}
	num := selectNamespace(nss)
	co.namespace = nss[num].Name
	return nil
}

func (co *CreateOptions) createServiceAccounts() error {
	saName := co.userName
	userServiceAccount, err := co.clientSet.CoreV1().ServiceAccounts(co.namespace).Get(context.TODO(), saName, metav1.GetOptions{})
	if err != nil {
		saObj := &coreV1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: saName,
			},
		}
		userServiceAccount, err = co.clientSet.CoreV1().ServiceAccounts(co.namespace).Create(context.TODO(), saObj, metav1.CreateOptions{})
		if err != nil {
			return err
		}
		printString(os.Stdout, "ServiceAccount")
		fmt.Printf(" : %s create success\n", userServiceAccount.Name)
	} else {
		printYellow(os.Stdout, "ServiceAccount")
		fmt.Printf(" : %s already exists\n", userServiceAccount.Name)
	}
	return nil
}

func (co *CreateOptions) selectClusterRole() error {
	clusterRoleList := []string{
		"view", "edit", "admin", "cluster-admin", "custom",
	}
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F63C {{ . | red }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U0001F638 Select:{{ . | green }}",
	}
	prompt := promptui.Select{
		Label:     "please select the cluster role of the user:",
		Items:     clusterRoleList,
		Templates: templates,
		Size:      uiSize,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	if clusterRoleList[i] == "custom" {
		customClusterRole := PromptUI("custom cluster role", "")
		fmt.Println(customClusterRole)
		_, err := co.clientSet.RbacV1().ClusterRoles().Get(context.TODO(), customClusterRole, metav1.GetOptions{})
		if err != nil {
			return err
		}
		co.role = customClusterRole
	} else {
		co.role = clusterRoleList[i]
	}
	return nil
}

func (co *CreateOptions) createRoleBinding() error {
	rb := &rbacV1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", co.userName, co.role),
			Namespace: co.namespace,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      co.userName,
				Namespace: co.namespace,
			},
		},
		RoleRef: rbacV1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     co.role,
		},
	}
	newRoleBinding, err := co.clientSet.RbacV1().RoleBindings(co.namespace).Create(context.TODO(), rb, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	printString(os.Stdout, "RoleBinding")
	fmt.Printf(" : %s create success\n", newRoleBinding.Name)
	return nil
}

func (co *CreateOptions) getToken() error {
	sa, err := co.clientSet.CoreV1().ServiceAccounts(co.namespace).Get(context.TODO(), co.userName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	secretName := sa.Secrets[0].Name
	secretToken, _ := co.clientSet.CoreV1().Secrets(co.namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	sEnc := base64.StdEncoding.EncodeToString(secretToken.Data["token"])
	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	if err != nil {
		return err
	}
	co.token = string(sDec)
	return nil
}

func (co *CreateOptions) putOutKubeConfig() *clientcmdapi.Config {
	coContext := co.config.Contexts[co.contextName]
	coCluster := co.config.Clusters[coContext.Cluster]
	coAuthInfo := clientcmdapi.NewAuthInfo()
	coAuthInfo.Token = co.token
	coContext.AuthInfo = co.userName

	newConfig := clientcmdapi.NewConfig()
	newConfig.Clusters[coContext.Cluster] = coCluster
	newConfig.AuthInfos[coContext.AuthInfo] = coAuthInfo
	newConfig.Contexts[co.userName] = coContext
	newConfig.CurrentContext = co.userName
	return newConfig
}

func createExample() string {
	return `
# Create new KubeConfig(experiment)
kubecm create
`
}
