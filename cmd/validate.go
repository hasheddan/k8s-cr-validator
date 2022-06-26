package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/moulick/k8s-cr-validator/validate"
)

var (
	crFiles    []string
	crFolders  []string
	crdFiles   []string
	crdFolders []string
	ignoreKind []string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate CR againt CRD",
	Long:  "Validate CR againt CRD",
	// Args: cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validate.Validate(crFiles, crFolders, crdFiles, crdFolders, ignoreKind, k8sVersion); err != nil {
			// red color for errors
			color.Set(color.FgHiRed, color.Bold)
			defer color.Unset()
			// Assert that colors will definitely be used if requested
			// if forceColor {
			// 	color.NoColor = false // TODO: make this work
			// }
			if merr, ok := err.(*multierror.Error); ok {
				fmt.Printf("Number of errors %v\n", merr.Len())
				for _, e := range merr.Errors {
					fmt.Println(e)
				}
				os.Exit(1)
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		color.Green("All Good")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// individual cr files
	validateCmd.Flags().StringSliceVarP(&crFiles, "cr-files", "c", []string{}, "Comma separated list of files containing Kubernetes CR(s) (can be specified multiple times)")
	// cr folders
	validateCmd.Flags().StringSliceVar(&crFolders, "cr-folders", []string{}, "Comma separated list of folders containing Kubernetes CR(s) (can be specified multiple times)")
	// individual crd files
	validateCmd.Flags().StringSliceVarP(&crdFiles, "crd-files", "d", []string{}, "Comma separated list of files containing Kubernetes CRD(s) (can be specified multiple times)")
	// crd folders
	validateCmd.Flags().StringSliceVar(&crdFolders, "crd-folders", []string{}, "Comma separated list of folders containing Kubernetes CRD(s) (can be specified multiple times)")
	validateCmd.Flags().StringSliceVar(&ignoreKind, "ignore-kinds", []string{}, "Comma separated list of Kinds to ignore (can be specified multiple times)")
}
