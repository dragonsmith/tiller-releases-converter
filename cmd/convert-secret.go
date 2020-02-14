package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var convertSecretCmd = &cobra.Command{
	Use:   "convert-secret ARGS",
	Short: "Convert Tiller Secrets releases to Configmaps.",
	Long: `Convert (tiller-releseases-converter convert) will actually create a new ConfigMap
for each Tiller-owned Secret.`,

	Run: func(cmd *cobra.Command, args []string) {
		tillerSecrets, err := listTillerSecrets()
		if err != nil {
			log.Fatalln(err)
		}

		for _, item := range tillerSecrets {
			fmt.Printf(" - [    ] %s", item.ObjectMeta.Name)

			err := createConfigMapFromSecret(item)
			if err != nil && apierrors.IsAlreadyExists(err) {
				fmt.Printf("\r - [%s] %s (target already exists)\n", color.RedString("FAIL"), item.ObjectMeta.Name)
			} else {
				fmt.Printf("\r - [ %s ] %s\n", color.GreenString("OK"), item.ObjectMeta.Name)
			}
		}
	},
}

func createConfigMapFromSecret(secret corev1.Secret) error {
	var valueAsString = string(secret.Data["release"])

	newConfigMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   secret.ObjectMeta.Name,
			Labels: secret.ObjectMeta.Labels,
		},
		Data: map[string]string{"release": valueAsString},
	}

	//Create a new Secret resource with data from an old ConfigMap
	_, err := clientset.CoreV1().ConfigMaps(destinationNameSpace).Create(&newConfigMap)

	return err
}
