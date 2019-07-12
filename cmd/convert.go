package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const (
	tillerDeploymentName string = "tiller-deploy"
)

var convertCmd = &cobra.Command{
	Use:   "convert ARGS",
	Short: "Convert Tiller ConfigMap releases to Secrets.",
	Long: `Convert (tiller-releseases-converter convert) will actually create a new Secret
for each Tiller-owned ConfigMap.`,

	Run: func(cmd *cobra.Command, args []string) {
		tillerConfigMaps, err := listTillerConfigMaps()
		if err != nil {
			log.Fatalln(err)
		}

		for _, item := range tillerConfigMaps {
			fmt.Printf(" - [    ] %s", item.ObjectMeta.Name)

			err := createSecretFromConfigMap(item)
			if err != nil && apierrors.IsAlreadyExists(err) {
				fmt.Printf("\r - [%s] %s (target already exists)\n", color.RedString("FAIL"), item.ObjectMeta.Name)
			} else {
				fmt.Printf("\r - [ %s ] %s\n", color.GreenString("OK"), item.ObjectMeta.Name)
			}
		}
	},
}

func createSecretFromConfigMap(configMap corev1.ConfigMap) error {
	newSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   configMap.ObjectMeta.Name,
			Labels: configMap.ObjectMeta.Labels,
		},
		Data: map[string][]byte{"release": []byte(configMap.Data["release"])},
	}

	// Create a new Secret resource with data from an old ConfigMap
	_, err := clientset.CoreV1().Secrets(destinationNameSpace).Create(&newSecret)

	return err
}
