package internal

import (
	"log"

	"github.com/MeteorsLiu/llpkgstore/actions"
	"github.com/spf13/cobra"
)

var postProcessingCmd = &cobra.Command{
	Use:   "post-processing",
	Short: "Verify a PR",
	Long:  ``,
	Run:   runPostProcessingCmd,
}

func runPostProcessingCmd(_ *cobra.Command, _ []string) {
	log.Println("run post processing")
	actions.NewDefaultClient().Release()
}

func init() {
	rootCmd.AddCommand(postProcessingCmd)
}
