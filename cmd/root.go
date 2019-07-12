package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeContext          string
	kubeConfig           string
	nameSpace            string
	destinationNameSpace string

	clientset *kubernetes.Clientset
)

func init() {
	cobra.OnInitialize(initConfig)

	convertCmd.LocalFlags().StringVarP(&destinationNameSpace, "destination-namespace", "d", "", "destination tiller namespace (default is set to namespace flag)")

	rootCmd.PersistentFlags().StringVar(&kubeContext, "context", "", "kube config context")
	rootCmd.PersistentFlags().StringVarP(&kubeConfig, "kubeconfig", "c", "", "config file (default is $HOME/.kube/config)")
	rootCmd.PersistentFlags().StringVarP(&nameSpace, "namespace", "n", "", "tiller namespace (default is kube-system)")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(secureTillerCmd)
}

var rootCmd = &cobra.Command{
	Use:   "tiller-releases-converter ARGS",
	Short: "A converter for Tiller's releases from ConfigMaps to Secrets",
	Long: `A converter for Tiller's releases from ConfigMaps to Secrets
to migrate a default Tiller installation to a more secure one,
which uses K8s secrets as its backend.`,
	Version: "0.1.1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func initConfig() {

	// Setting defaults
	if kubeConfig == "" {
		kubeConfig = filepath.Join(homeDir(), ".kube", "config")
	}

	if nameSpace == "" {
		nameSpace = "kube-system"
	}

	if destinationNameSpace == "" {
		destinationNameSpace = nameSpace
	}

	config, err := getKubeConfig()
	if err != nil {
		log.Fatalln("Cannot assemble client configuration from Kube config:\n\n", err)
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("Cannot create a Kubernetes connection:\n\n", err)
	}
}

// GetConfig returns a Kubernetes client config for a given context.
func getKubeConfig() (*restclient.Config, error) {
	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig}
	overrides := &clientcmd.ConfigOverrides{}

	if kubeContext != "" {
		overrides.CurrentContext = kubeContext
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).ClientConfig()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
