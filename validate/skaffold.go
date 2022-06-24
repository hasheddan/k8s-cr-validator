package validate

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
)

// Validate returns an error
// if the cr(s) cannot be validated againt the given crd(s)
// or else if crd is not given and does not conform to jsonSchema(s) for known kubernetes native objects
func Validate(crFiles []string, crFolders []string, crdFiles []string, crdFolders []string, ignoreKind []string) error {
	var merr *multierror.Error

	// collect paths of all CRDs
	crdPaths := aggregateFiles(crdFiles, crdFolders)
	// get the generated validators
	validators, err := generateValidators(crdPaths)
	if err != nil {
		merr = multierror.Append(merr, fmt.Errorf("failed to generate validators: %w", err))
	}

	// collect paths of all CRs
	crPaths := aggregateFiles(crFiles, crFolders)
	// fmt.Printf("List of files to Validate %v\n", crPaths)
	crList, err := readCR(crPaths, ignoreKind)
	if err != nil {
		merr = multierror.Append(merr, fmt.Errorf("failed to read CR: %w", err))
	}

	// First check if validator is available, else it's possible it's a native oject, hence try validating by kubeConform
	for _, cr := range crList {
		v, ok := validators[cr.GroupVersionKind()]
		// didn't find validator, try kubeConform
		if !ok {
			fmt.Printf(
				"checking %s (%s) againt kubeconform... ",
				cr.GetName(),
				cr.GroupVersionKind(),
			)
			err = kubeConform(cr)
			if err != nil {
				merr = multierror.Append(merr, fmt.Errorf("failed to validate CR: %w", err))
				continue
			}
			continue
		}

		// now validate againt our matched validator
		result := v.Validate(cr)
		if len(result.Errors) > 0 || len(result.Warnings) > 0 {
			for _, e := range result.Errors {
				merr = multierror.Append(merr, fmt.Errorf(
					"validation Error %s, %s: %w",
					cr.GroupVersionKind(),
					cr.GetName(),
					e,
				))
			}
			for _, e := range result.Warnings {
				merr = multierror.Append(merr, fmt.Errorf("validation Warning %s, %s: %w",
					cr.GroupVersionKind().String(),
					cr.GetName(),
					e,
				))
			}
		} else {
			fmt.Printf(
				"CR %s (%s) is %s \n",
				cr.GetName(),
				cr.GroupVersionKind(),
				color.GreenString("Valid ✅"),
			)
		}
	}
	return merr.ErrorOrNil()
}
