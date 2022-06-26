# k8s-cr-validator

## Badges

[![Build status](https://img.shields.io/github/workflow/status/moulick/k8s-cr-validator/goreleaser?style=for-the-badge)](https://github.com/moulick/k8s-cr-validator/actions?workflow=goreleaser)
[![Release](https://img.shields.io/github/v/release/moulick/k8s-cr-validator?style=for-the-badge)](https://github.com/moulick/k8s-cr-validator/releases/latest)
[![Software License](https://img.shields.io/github/license/moulick/k8s-cr-validator?style=for-the-badge)](/LICENSE.md)
[![Go Report card](https://goreportcard.com/badge/github.com/moulick/k8s-cr-validator?style=for-the-badge)](https://goreportcard.com/report/github.com/moulick/k8s-cr-validator)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/moulick/k8s-cr-validator)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)

k8s-cr-validator is a [Kubernetes custom resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) validation tool.
It generates validators for the custom resource definitions and validates a given custom resource object against it. The magic is that k8s-cr-validator uses the same code to generate validators as the Kubernetes API Server

For resources it does not find a CRD, it will try to validate it against [KubeConform](https://github.com/yannh/kubeconform). This is mainly intended for native Kubernetes resources like namespaces or deployments etc.

### Comparison to KubeConform

KubeConform requires converting your CRDs to JSON schema and storing the JSON schemas in a local folder. Across the internet, you'll find CRDs as YAML files, be it from helm charts or standalone files. This tool takes out the manual step of generating and storing the CRDs as JSON Schemas.

### Usage

```bash
Validate CR againt CRD

Usage:
  k8s-cr-validator validate [flags]

Flags:
  -c, --cr-files strings       Comma separated list of files containing Kubernetes CR(s) (can be specified multiple times)
      --cr-folders strings     Comma separated list of folders containing Kubernetes CR(s) (can be specified multiple times)
  -d, --crd-files strings      Comma separated list of files containing Kubernetes CRD(s) (can be specified multiple times)
      --crd-folders strings    Comma separated list of folders containing Kubernetes CRD(s) (can be specified multiple times)
  -h, --help                   help for validate
      --ignore-kinds strings   Comma separated list of Kinds to ignore (can be specified multiple times)

Global Flags:
      --kubernetes-version string   Version of Kubernetes to validate against, e.g: 1.20.12 (default "master")
```

### Usage examples

* Validating a single, invalid file against a single CRDs

```bash
$ k8s-cr-validator validate --crd-files ./examples/crds/crossplane.yaml --cr-files ./examples/cr/cr-invalid.yaml
Found Validator for apiextensions.crossplane.io/v1, Kind=CompositeResourceDefinition
Found Validator for apiextensions.crossplane.io/v1beta1, Kind=CompositeResourceDefinition
Found Validator for apiextensions.crossplane.io/v1, Kind=Composition
Found Validator for apiextensions.crossplane.io/v1beta1, Kind=Composition
Found Validator for apiextensions.crossplane.io/v1alpha1, Kind=CompositionRevision
Found Validator for meta.pkg.crossplane.io/v1, Kind=Configuration
Found Validator for meta.pkg.crossplane.io/v1alpha1, Kind=Configuration
Found Validator for meta.pkg.crossplane.io/v1, Kind=Provider
Found Validator for meta.pkg.crossplane.io/v1alpha1, Kind=Provider
Found Validator for pkg.crossplane.io/v1, Kind=Configuration
Found Validator for pkg.crossplane.io/v1beta1, Kind=Configuration
Found Validator for pkg.crossplane.io/v1, Kind=ConfigurationRevision
Found Validator for pkg.crossplane.io/v1beta1, Kind=ConfigurationRevision
Found Validator for pkg.crossplane.io/v1, Kind=Provider
Found Validator for pkg.crossplane.io/v1beta1, Kind=Provider
Found Validator for pkg.crossplane.io/v1, Kind=ProviderRevision
Found Validator for pkg.crossplane.io/v1beta1, Kind=ProviderRevision
Found Validator for pkg.crossplane.io/v1alpha1, Kind=ControllerConfig
Found Validator for pkg.crossplane.io/v1alpha1, Kind=Lock
Found Validator for pkg.crossplane.io/v1beta1, Kind=Lock
Loaded example1 of type apiextensions.crossplane.io/v1, Kind=Composition
Loaded example2 of type apiextensions.crossplane.io/v1, Kind=Composition
Number of errors 4
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example1: spec.resources.connectionDetails.fromConnectionSecretKey in body must be of type string: "number"
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example1: spec.resources.patches.transforms.type in body is required
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example2: spec.resources.connectionDetails.fromConnectionSecretKey in body must be of type string: "number"
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example2: spec.resources.patches.transforms.type in body is required
exit status 1
$ echo $?
1
```

* Validating a single valid file against a single CRD

```bash
$ k8s-cr-validator validate --crd-files ./examples/crds/external-secrets.yaml --cr-files ./examples/cr/externalsecret.yaml
Found Validator for kubernetes-client.io/v1, Kind=ExternalSecret
Loaded aws-secretsmanager of type kubernetes-client.io/v1, Kind=ExternalSecret
CR aws-secretsmanager (kubernetes-client.io/v1, Kind=ExternalSecret) is Valid ✅
All Good
$ echo $?
0
```

* Validating a folder of custom resources against a folder of CRDs

```bash
$ k8s-cr-validator validate --crd-folders ./examples/crds --cr-folders ./examples/cr
Found Validator for apiextensions.crossplane.io/v1, Kind=CompositeResourceDefinition
Found Validator for apiextensions.crossplane.io/v1beta1, Kind=CompositeResourceDefinition
Found Validator for apiextensions.crossplane.io/v1, Kind=Composition
Found Validator for apiextensions.crossplane.io/v1beta1, Kind=Composition
Found Validator for apiextensions.crossplane.io/v1alpha1, Kind=CompositionRevision
Found Validator for meta.pkg.crossplane.io/v1, Kind=Configuration
Found Validator for meta.pkg.crossplane.io/v1alpha1, Kind=Configuration
Found Validator for meta.pkg.crossplane.io/v1, Kind=Provider
Found Validator for meta.pkg.crossplane.io/v1alpha1, Kind=Provider
Found Validator for pkg.crossplane.io/v1, Kind=Configuration
Found Validator for pkg.crossplane.io/v1beta1, Kind=Configuration
Found Validator for pkg.crossplane.io/v1, Kind=ConfigurationRevision
Found Validator for pkg.crossplane.io/v1beta1, Kind=ConfigurationRevision
Found Validator for pkg.crossplane.io/v1, Kind=Provider
Found Validator for pkg.crossplane.io/v1beta1, Kind=Provider
Found Validator for pkg.crossplane.io/v1, Kind=ProviderRevision
Found Validator for pkg.crossplane.io/v1beta1, Kind=ProviderRevision
Found Validator for pkg.crossplane.io/v1alpha1, Kind=ControllerConfig
Found Validator for pkg.crossplane.io/v1alpha1, Kind=Lock
Found Validator for pkg.crossplane.io/v1beta1, Kind=Lock
Found Validator for kubernetes-client.io/v1, Kind=ExternalSecret
Loaded example1 of type apiextensions.crossplane.io/v1, Kind=Composition
Loaded example2 of type apiextensions.crossplane.io/v1, Kind=Composition
Loaded aws-secretsmanager of type kubernetes-client.io/v1, Kind=ExternalSecret
Loaded foo of type /v1, Kind=Namespace
Loaded bar of type /v1, Kind=Namespace
CR aws-secretsmanager (kubernetes-client.io/v1, Kind=ExternalSecret) is Valid ✅
checking foo (/v1, Kind=Namespace) againt kubeconform... Valid ✅
checking bar (/v1, Kind=Namespace) againt kubeconform... Validate Failed ❌
Number of errors 5
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example1: spec.resources.patches.transforms.type in body is required
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example1: spec.resources.connectionDetails.fromConnectionSecretKey in body must be of type string: "number"
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example2: spec.resources.patches.transforms.type in body is required
validation Error apiextensions.crossplane.io/v1, Kind=Composition, example2: spec.resources.connectionDetails.fromConnectionSecretKey in body must be of type string: "number"
failed to validate CR: failed to validate bar (/v1, Kind=Namespace): For field metadata: Additional property type is not allowed
exit status 1
```



## Credits

* This was forked from https://github.com/hasheddan/k8s-cr-validator and is based on the excellent blog at http://danielmangum.com/posts/how-kubernetes-validates-custom-resources
* This tool also utilizes [KubeConform](https://github.com/yannh/kubeconform)
