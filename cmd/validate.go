package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/hasheddan/k8s-cr-validator/validate"
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
		if err := validate.Validate(crFiles, crFolders, crdFiles, crdFolders, ignoreKind); err != nil {
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
	validateCmd.Flags().StringSliceVar(&crFiles, "cr-files", []string{}, "List of files to validate. Repeat the flag for multiple files")
	// cr folders
	validateCmd.Flags().StringSliceVar(&crFolders, "cr-folders", []string{}, "List of folders containing files to validate. Repeat the flag for multiple folders")

	// individual crd files
	validateCmd.Flags().StringSliceVar(&crdFiles, "crd-files", []string{}, "List of files containing CRD(s). Repeat the flag for multiple files")
	// crd folders
	validateCmd.Flags().StringSliceVar(&crdFolders, "crd-folders", []string{}, "List of folders containing CRD(s). Repeat the flag for multiple folders")

	validateCmd.Flags().StringSliceVar(&ignoreKind, "ignore-kinds", []string{}, "List of Kinds to ignore. Repeat the flag for multiple kinds")

}
