package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup ARGS",
	Short: "Delete Tiller ConfigMap releases.",
	Long:  "Cleanup a cluster from old ConfigMap-based Tiller releases.",

	Run: func(cmd *cobra.Command, args []string) {
		tillerConfigMaps, err := listTillerConfigMaps()
		if err != nil {
			log.Fatalln(err)
		}

		for _, item := range tillerConfigMaps {
			fmt.Printf("Deleting %s: ", item.ObjectMeta.Name)
			if err = clientset.CoreV1().ConfigMaps(nameSpace).Delete(item.ObjectMeta.Name, &metav1.DeleteOptions{}); err != nil {
				color.Red("FAIL")
				log.Fatalln(err)
			}

			color.Green("OK")
		}
	},
}
