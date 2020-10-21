# Clowder :cat: - Insights Platform Operator

An operator to deploy and operate cloud.redhat.com resources for Openshift.

<img src="clowder.svg" width="150" alt="Clowder Logo">

## Overview

The Clowder operator takes much of the heavy lifting out of creating and 
maintaining applications on the Clouddot platform. Instead of an app developer
being responsible for authoring multiple resources and combining them into a
single k8s template, the Clowder app defines a simple `ClowdApp` resource
which not only defines the pods for the application, but also requests certain
key resources from the environment, such as Object Storage Buckets, a Database,
Kafka Topics, or an In Memory Database.

A `ClowdEnvironment` resource is used to define how key resources, such as
Kafka Topics and Databases are deployed. Using different providers, a single
`ClowdEnvironment` can radically alter the way in which resources are
provisioned. For example, with regards to Kafka Topics, setting the `provider`
to `local` will instruct Clowder to deploy a local Zookeeper/Kafka pod
and create topics inside it, but if the provider were set to `operator`, then
Clowder would instead drop a KafkaTopic custom resource ready for the Kafka
Strimzi operator to pick up and create topics.

The diagram below shows how the two Clowder resources are used to create all
other k8s resources.

![Clowder Flow](images/clowder-flow.svg "Clowder Flow")

Once these custom resources have been created and deployed to the k8s
environment, the operator will create a secret with all necessary configuration
data and expose it to the pods by mounting the JSON document in the app 
container. In this way, instead of an app configuring itself, the app is
configured instead by Clowder.

This has the advantage of creating consistency across deployments, whether
they are development, testing or production. It also creates a simple interface
for developers to onboard, producing a more streamlined developer experience.

The `ClowdApp` resource does not change when deploying into environments
configured with different `ClowdEnvironment` resources. Though the underlying
environmental resources, object storage, kafka, etc, may be provided through
different implementations, the configuration that is presented to the pod
remains consistent.

If the application is written in Python or Go, there is a client available
which further simplifies the process of obtaining configuration and offers
several helpers for accessing some of the more complex structures.

The diagram below describes how the application accesses the configuration.

!["Clowder Client"](images/clowder-new.svg "Clowder Client")

## Design

[Design docs](https://github.com/RedHatInsights/clowder/tree/master/docs/)

## Dependencies

- [Operator SDK](https://github.com/operator-framework/operator-sdk/releases)
- [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder/releases)
- [kustomize](https://github.com/kubernetes-sigs/kustomize/releases)
- Either Codeready Containers or a remote cluster where you have access to
  create CRDs.

## Running

- `make install` will deploy the CRDs to the cluster configured in your kubeconfig.
- `make run` will build the binary and locally run the binary, connecting the
  manager to the Openshift cluster configured in your kubeconfig.
- `make deploy` will try to run the manager in the cluster configured in your
  kubeconfig.  You likely need to push the image to an image stream local to
  your target namespace.
- `make genconfig` (optionally) needs to be run if the specification for the config
  has been altered.

## Testing

The tests rely on the test environment set up by controller-runtime.  This
enables the operator to get initialized against a control pane just like it
would against a real Openshift cluster.

While the tests do not rely on any additional testing frameworks (e.g. Ginkgo),
you do need to download
[kubebuilder](https://github.com/kubernetes-sigs/kubebuilder/releases) in order
to set up the control plane used
by the controller-runtime test environment.

Run the tests:

```
$ KUBEBUILDER_ASSETS=$PWD/kubebuilder go test ./controllers/...
ok      cloud.redhat.com/clowder/v2/controllers    9.626s
```

## Kuttl Tests

```
$ KUBEBUILDER_ASSETS=$HOME/workspace/kubebuilder kubectl kuttl test --config test-config.yaml 
...
    kuttl: harness.go:326: run tests finished
    kuttl: harness.go:401: killing process "./bin/manager"
    kuttl: harness.go:411: bg process: "./bin/manager" exit code -1
    kuttl: harness.go:427: tearing down mock control plane
--- PASS: kuttl (6.27s)
    --- PASS: kuttl/harness (0.00s)
        --- PASS: kuttl/harness/test_local_kafka (1.08s)
PASS
```
