package cmd

import (
	"fmt"
	"kubedev/pkg/build/image"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build images of the kubernetes components",
	Long:  "Build images of the kubernetes components",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		err := BuildComponents(args)
		if err != nil {
			Fatal(err.Error(), DefaultErrorExitCode)
		}
	},
}

func BuildComponents(args []string) error {
	return image.BuildImages()
}
