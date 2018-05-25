package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/util/retry"
)

var secureTillerCmd = &cobra.Command{
	Use:   "secure-tiller ARGS",
	Short: "Patch Tiller Deployment to use secrets storage.",
	Long: `This command patches in-cluster Tiller deployment adding
"--storage=secrets" attribute to tiller cmd.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating Tiller Deployment...")
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			deploymentsClient := clientset.AppsV1().Deployments(nameSpace)
			result, getErr := deploymentsClient.Get(tillerDeploymentName, metav1.GetOptions{})
			if getErr != nil {
				log.Println("Failed to get latest version of Deployment...")
				return getErr
			}

			result.Spec.Template.Spec.Containers[0].Command = []string{"/tiller", "--storage=secret"}
			_, updateErr := deploymentsClient.Update(result)
			return updateErr
		})

		if retryErr != nil {
			log.Fatalf("Update failed: %v", retryErr)
		}
		fmt.Println("Tiller Deployment was updated successfully!")
	},
}
