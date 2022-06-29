> :warning: **This repo is archived! Go to**: https://github.com/cisco-open/cluster-registry-controller :warning:


# Cluster Registry

This repository defines a lightweight Kubernetes Custom Resource Definition API
for defining a list of clusters and associated metadata in a K8s environment.

## What is it?

The cluster registry helps you keep track of and perform operations on your
clusters within one administrative domain. This repository doesn't contain
an actual implementation of a controller that uses the API.

## Defined CRDs

1. `Cluster`: defines a Kubernetes cluster.
2. `ResourceSyncRule`: defines a sync rule based on which Kubernetes resources are synced across clusters.
3. `ClusterFeature`: defines a feature name, which can be used by a `ResourceSyncRule` resource to define which clusters
to sync from a given Kubernetes resource.
