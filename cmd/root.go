package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// forceColor bool
	k8sVersion string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k8s-cr-validator",
	Short: "Validate K8s CR against CRD",
	Long:  `k8s-cr-validator is a tool that is used to validate a Kubernetes Custom Resource againt a Custom Resource Definition`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// rootCmd.PersistentFlags().BoolVarP(&forceColor, "force-color", "", false, "Force colored output even if stdout is not a TTY")
	rootCmd.PersistentFlags().StringVar(&k8sVersion, "kubernetes-version", "master", "Version of Kubernetes to validate against, e.g: 1.20.12")

}
