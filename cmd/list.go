package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kblabels "k8s.io/apimachinery/pkg/labels"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var listCmd = &cobra.Command{
	Use:   "list ARGS",
	Short: "List Tiller ConfigMap releases.",
	Long: `List (tiller-releseases-converter list) will list all Tiller's ConfigMaps
in the designated namespace.`,

	Run: func(cmd *cobra.Command, args []string) {
		tillerConfigMaps, err := listTillerConfigMaps()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("I've found these Tiller's ConfigMap releases for you:")
		fmt.Println()
		for _, item := range tillerConfigMaps {
			fmt.Println(item.ObjectMeta.Name)
		}
	},
}

func listTillerConfigMaps() ([]corev1.ConfigMap, error) {
	lsel := kblabels.Set{"OWNER": "TILLER"}.AsSelector()
	opts := metav1.ListOptions{LabelSelector: lsel.String()}

	configMaps, err := clientset.CoreV1().ConfigMaps(nameSpace).List(opts)
	if err != nil {
		return nil, err
	}

	return configMaps.Items, nil
}

func listTillerSecrets() ([]corev1.Secret, error) {
	lsel := kblabels.Set{"OWNER": "TILLER"}.AsSelector()
	opts := metav1.ListOptions{LabelSelector: lsel.String()}

	secrets, err := clientset.CoreV1().Secrets(nameSpace).List(opts)
	if err != nil {
		return nil, err
	}

	return secrets.Items, nil
}
