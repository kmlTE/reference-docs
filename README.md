# Kubernetes Reference Docs

Tools to build reference documentation for Kubernetes APIs and CLIs.

See [Generating Reference Documentation for kubectl Commands](https://kubernetes.io/docs/contribute/generate-ref-docs/kubectl/) for information on how to generate kubectl reference docs.

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).
Instructions on how to contribute can be found [here](CONTRIBUTING.md).

You can reach the maintainers of this project at:

- [Slack channel](https://kubernetes.slack.com/messages/sig-docs)
- [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-sig-docs)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

## TOSCA Generation

From the Kubernetes API specs, the TOSCA definitions will be generated.

Prerequisites:
- Python v3.7.x
- Git
- Golang version 1.13+
- Pip used to install PyYAML
- PyYAML v5.1.2
- make
- gcc compiler/linker
- Docker (Required only for kubectl command reference)

To prepare the workspace and set `$GOPATH`:

```bash
# create workspace
mkdir -p $PWD/kube-doc-model
export GOPATH=$PWD/kube-doc-model

# making website destinaton path
mkdir -p $GOPATH/src/github.com/dummy/website

# dependencies
go get -u github.com/kmlTE/reference-docs
go get -u github.com/go-openapi/loads
go get -u github.com/go-openapi/spec
git clone https://github.com/kubernetes/kubernetes $GOPATH/src/k8s.io/kubernetes
```

To generate the Kubernetes TOSCA data types and node types, the environment variables must be set and `make` should be run:

```bash
export K8S_WEBROOT=$GOPATH/src/github.com/dummy/website # <web-base>
export K8S_ROOT=$GOPATH/src/k8s.io/kubernetes # <k8s-base>
export RDOCS_ROOT=$GOPATH/src/github.com/kmlTE/reference-docs # <rdocs-base>
export K8S_RELEASE=1.18.12

make api
```

The output YAML file can then be found in `/tmp/kubernetes/kubernetes_definitions.yaml`. It contains TOSCA definitions for the following Kubernetes Kinds: *Deployment, ServiceAccount, ClusterRole, ClusterRoleBinding, Namespace, DaemonSet*, and is based on v1_18 spec. To add other definitions, modify `included_objects` in configuration file `gen-apidocs/config/v1_18/config.yaml`.

To use another spec version, change `$K8S_RELEASE` and add needed definitions in respective configuration file, e.g. in `gen-apidocs/config/v1_19/config.yaml`:
```YAML
included_objects:
  - "Deployment"
  - "Namespace"
  - "DaemonSet"
...
```