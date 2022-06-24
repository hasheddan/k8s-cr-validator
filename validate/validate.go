package validate

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/yannh/kubeconform/pkg/resource"
	"github.com/yannh/kubeconform/pkg/validator"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/apiserver/validation"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiMachYaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/kube-openapi/pkg/validation/validate"
	"sigs.k8s.io/yaml"
)

// aggregateFiles returns a list of aggregated paths for given list of file and folders
func aggregateFiles(files []string, directories []string) []string {
	var paths []string
	// for files simply add them to the list
	for _, file := range files {
		paths = append(paths, file)
	}

	// for directories, recursively add all files in the directory
	for _, directory := range directories {
		err := fs.WalkDir(os.DirFS(directory), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			// only consider files with .yaml or .yml extension
			if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
				paths = append(paths, filepath.Join(directory, path))
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	return paths
}

// generateValidators returns a map of validate.SchemaValidator for all the given CRD files
func generateValidators(crds []string) (map[schema.GroupVersionKind]*validate.SchemaValidator, error) {
	validators := map[schema.GroupVersionKind]*validate.SchemaValidator{}

	for _, file := range crds {
		f, err := os.Open(file)
		if err != nil {
			f.Close()
			return nil, err
		}

		yr := apiMachYaml.NewYAMLReader(bufio.NewReader(f))
		// loop to get all documents from a single file
		for {
			b, err2 := yr.Read()
			if err2 != nil && err2 != io.EOF {
				return nil, err2
			}
			if err2 == io.EOF {
				break
			}
			if len(b) == 0 {
				continue
			}
			crd := &extv1.CustomResourceDefinition{}
			if err = yaml.Unmarshal(b, crd); err != nil {
				f.Close()
				return nil, fmt.Errorf("error parsing CRD %s: %v", b, err)
			}

			// convert the CRD to the internal representation
			internal := &apiextensions.CustomResourceDefinition{}
			if err = extv1.Convert_v1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition(crd, internal, nil); err != nil {
				f.Close()
				return nil, fmt.Errorf("error converting V1 CRD to APIExtensions CRD  %s: %v", b, err)
			}
			// loop over versions stored in the CRD
			for _, ver := range internal.Spec.Versions {
				gvk := schema.GroupVersionKind{
					Group:   internal.Spec.Group,
					Version: ver.Name,
					Kind:    internal.Spec.Names.Kind,
				}

				var schemaValidator *validate.SchemaValidator
				schemaValidator, _, err = validation.NewSchemaValidator(ver.Schema)
				if err != nil {
					f.Close()
					return nil, fmt.Errorf("error creating schema validator for CRD %s in file %s: %v", gvk, file, err)
				}
				if internal.Spec.Validation != nil {
					schemaValidator, _, err = validation.NewSchemaValidator(internal.Spec.Validation)
					if err != nil {
						f.Close()
						return nil, fmt.Errorf("error parsing validation schema for CRD %s in file %s: %v", gvk, file, err)
					}
				}

				// check if there is any duplicate validator/crd
				// prevent unnecessary frustration for the user incase there is a duplicate
				if _, exists := validators[gvk]; exists {
					f.Close()
					return nil, fmt.Errorf("duplicate CRD found for %s", gvk)
				}
				validators[gvk] = schemaValidator
				color.Green("Found Validator for %s", gvk)
			}
		}
		err = f.Close()
		if err != nil {
			return nil, err
		}
	}
	return validators, nil
}

// readCR returns unmarshalled unstructured.Unstructured for the given CR files
func readCR(crs []string) ([]*unstructured.Unstructured, error) {
	var crList []*unstructured.Unstructured

	for _, crPath := range crs {
		if _, err := os.Stat(crPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("ERROR: %s does not exist, Error: %w", crPath, err)
		}

		f, err := os.Open(crPath)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("ERROR: opening %s: %w", crPath, err)
		}

		yr := apiMachYaml.NewYAMLReader(bufio.NewReader(f))
		// loop to get all documents from a single file
		for {
			b, err2 := yr.Read()
			if err2 != nil && err2 != io.EOF {
				return nil, err2
			}
			if err2 == io.EOF {
				break
			}
			if len(b) == 0 {
				continue
			}

			obj := &unstructured.Unstructured{}
			if err = yaml.Unmarshal(b, obj); err != nil {
				f.Close()
				return nil, fmt.Errorf("error decoding yaml %s: %w", b, err)
			}
			crList = append(crList, obj)
			color.Green("Loaded %s of type %s", obj.GetName(), obj.GroupVersionKind())
		}
		err = f.Close()
		if err != nil {
			return nil, fmt.Errorf("error closing file %s", crPath)
		}
	}
	return crList, nil
}

// kubeConform returns error if the given CR cannot be validated against by kubeConform
func kubeConform(cr *unstructured.Unstructured) error {

	// Create a validator
	conformValidator, err := validator.New(nil, validator.Opts{
		Strict:            true,
		KubernetesVersion: "1.20.12"})
	if err != nil {
		return fmt.Errorf("failed initializing validator: %w", err)
	}
	// convert cr to []byte
	obj, _ := cr.MarshalJSON()
	// validate the cr
	result := conformValidator.ValidateResource(resource.Resource{Bytes: obj})
	if result.Status != validator.Valid {
		color.Red("Validate Failed ❌\n")
		return fmt.Errorf(
			"failed to validate %s (%s): %w",
			cr.GetName(),
			cr.GroupVersionKind(),
			result.Err,
		)
	}
	color.Green("Valid ✅\n")
	return nil
}
