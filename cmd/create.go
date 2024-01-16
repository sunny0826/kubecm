package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	certificatesv1 "k8s.io/api/certificates/v1"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"os"
	"time"
)

// CreateCommand clean command struct
type CreateCommand struct {
	BaseCommand
}

type CreateOptions struct {
	config      *clientcmdapi.Config
	clientSet   kubernetes.Interface
	role        string
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
	//ce.command.DisableFlagsInUseLine = true
	ce.command.Flags().String("user", "", "user name for kubeconfig")
	ce.command.Flags().StringP("namespace", "n", "", "namespace for user")
	ce.command.Flags().String("cluster-role", "", "cluster role for user")
	ce.command.Flags().String("context-name", "", "context name for kubeconfig")
	ce.command.Flags().Bool("print-clean-up", false, "print clean up command")
}

func (ce *CreateCommand) runCreate(cmd *cobra.Command, args []string) error {
	userName, _ := ce.command.Flags().GetString("user")
	namespace, _ := ce.command.Flags().GetString("namespace")
	clusterRole, _ := ce.command.Flags().GetString("cluster-role")
	contextName, _ := ce.command.Flags().GetString("context-name")
	clean, _ := ce.command.Flags().GetBool("print-clean-up")

	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	if userName == "" {
		userName = PromptUI("user name", "")
	}
	co := CreateOptions{
		config:   config,
		userName: userName,
	}
	if contextName == "" {
		err = co.chooseContext()
		if err != nil {
			return err
		}
	} else {
		co.contextName = contextName
		co.config.CurrentContext = contextName
		set, err := GetClientSet(cfgFile)
		if err != nil {
			return err
		}
		co.clientSet = set
	}
	if namespace == "" {
		err = co.chooseNamespace()
		if err != nil {
			return err
		}
	} else {
		co.namespace = namespace
	}

	// create CSR
	_, privateKey, err := co.createCSR()
	if err != nil {
		return err
	}

	// approve CSR
	err = co.approveCSR()
	if err != nil {
		return err
	}

	if clusterRole == "" {
		// select ClusterRole
		err = co.selectClusterRole()
		if err != nil {
			return err
		}
	} else {
		co.role = clusterRole
	}

	// create RoleBinding
	err = co.createRoleBinding()
	if err != nil {
		return err
	}

	if clean {
		printCleanCmd(co.userName, co.role)
	}

	// create new kubeconfig
	return co.createKubeConfig(privateKey)
}

// createCSR create CSR
func (co *CreateOptions) createCSR() ([]byte, *rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   co.userName,
			Organization: []string{"kubecm"},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
		PublicKeyAlgorithm: x509.RSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return nil, nil, err
	}

	pemCSR := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

	csr := &certificatesv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: co.userName,
		},
		Spec: certificatesv1.CertificateSigningRequestSpec{
			Request:    pemCSR,
			Usages:     []certificatesv1.KeyUsage{certificatesv1.UsageDigitalSignature, certificatesv1.UsageKeyEncipherment, certificatesv1.UsageClientAuth},
			SignerName: certificatesv1.KubeAPIServerClientSignerName,
		},
	}

	csr, err = co.clientSet.CertificatesV1().CertificateSigningRequests().Create(context.TODO(), csr, metav1.CreateOptions{})
	if err != nil {
		return nil, nil, err
	}
	printString(os.Stdout, "CSR: "+csr.Name+" create success\n")
	return pemCSR, privateKey, err
}

// approveCSR approve CSR
func (co *CreateOptions) approveCSR() error {
	csr, err := co.clientSet.CertificatesV1().CertificateSigningRequests().Get(context.TODO(), co.userName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// ensure CSR is not approved
	for _, condition := range csr.Status.Conditions {
		if condition.Type == certificatesv1.CertificateApproved {
			printString(os.Stdout, "CSR: "+csr.Name+" has been approved\n")
			return nil
		}
	}

	// update to reflect approval status
	approvalCondition := certificatesv1.CertificateSigningRequestCondition{
		Type:           certificatesv1.CertificateApproved,
		Status:         coreV1.ConditionTrue, // Set Status == True
		Reason:         "ApprovedByAdmin",
		Message:        "This CSR was approved by admin.",
		LastUpdateTime: metav1.Now(),
	}

	csr.Status.Conditions = append(csr.Status.Conditions, approvalCondition)

	_, err = co.clientSet.CertificatesV1().CertificateSigningRequests().UpdateApproval(context.TODO(), co.userName, csr, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	printString(os.Stdout, "CSR: "+csr.Name+" has been approved\n")
	return err
}

// createKubeConfig create kubeconfig
func (co *CreateOptions) createKubeConfig(privateKey *rsa.PrivateKey) error {
	var csr *certificatesv1.CertificateSigningRequest
	var err error
	for i := 0; i < 10; i++ { // Retry up to 3 times
		csr, err = co.clientSet.CertificatesV1().CertificateSigningRequests().Get(context.TODO(), co.userName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if len(csr.Status.Certificate) != 0 {
			break
		}

		fmt.Printf("Waiting for CSR to be signed...  %v\n", i)
		// Sleep for a second before retrying
		time.Sleep(1 * time.Second)
	}

	certData := csr.Status.Certificate
	if len(certData) == 0 {
		return fmt.Errorf("certificate data is empty")
	}

	cluster := co.config.Clusters[co.contextName]
	if cluster == nil {
		return fmt.Errorf("cluster configuration not found")
	}

	newKubeConfig := clientcmdapi.NewConfig()
	newKubeConfig.Clusters[co.contextName] = &clientcmdapi.Cluster{
		Server:                   cluster.Server,
		CertificateAuthorityData: cluster.CertificateAuthorityData,
	}
	newKubeConfig.AuthInfos[co.userName] = &clientcmdapi.AuthInfo{
		ClientCertificateData: certData,
		ClientKeyData:         pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}),
	}
	newKubeConfig.Contexts[co.userName] = &clientcmdapi.Context{
		Cluster:  co.contextName,
		AuthInfo: co.userName,
	}
	newKubeConfig.CurrentContext = co.userName

	// write to file
	err = clientcmd.WriteToFile(*newKubeConfig, co.userName+"-kubeconfig.yaml")
	if err != nil {
		return err
	}
	printString(os.Stdout, "kubeconfig: "+co.userName+" create success\n")
	return err
}

// chooseContext choose context
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

// chooseNamespace choose namespace
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

// selectClusterRole select cluster role
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

// createRoleBinding create role binding
func (co *CreateOptions) createRoleBinding() error {
	rb := &rbacV1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", co.userName, co.role),
			Namespace: co.namespace,
		},
		Subjects: []rbacV1.Subject{
			{
				Kind:     "User",
				Name:     co.userName,
				APIGroup: "rbac.authorization.k8s.io",
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

func printCleanCmd(user, role string) {
	fmt.Print(`
# Clean up commands
kubectl delete certificatesigningrequests.certificates.k8s.io ` + user + `
kubectl delete rolebinding ` + user + `-` + role + `

`)
}

func createExample() string {
	return `
# Create new KubeConfig(experiment)
kubecm create
# Create new KubeConfig(experiment) with flags
kubecm create --user test --namespace default --cluster-role view --context-name kind-kind
`
}
