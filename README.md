# Sawtooth Kubernetes Operator

This is an initial test for the Dynamic Sawtooth P2P dynamic capability with seeds
being added on demand, 

## Installing

To install service accounts, roles and CRDs, follow:

```
$ kubectl create -f deploy/service_account.yaml

$ kubectl create -f deploy/role.yaml
$ kubectl create -f deploy/role_binding.yaml

# Install CRDs
$ kubectl create -f deploy/crds/sawtooth.hyperledger.org_sawtooths_crd.yaml
```

## Creating the Custom Resource

To start a new cluster with 5 nodes.

```
$ kubeclt apply -f crds/sawtooth.hyperledger.org_v1alpha1_sawtooth_cr.yaml

apiVersion: sawtooth.hyperledger.org/v1alpha1
kind: Sawtooth
metadata:
  name: sawtooth-1
spec:
  nodes: 5
  version: 1.2.3
  consensus: dev
```

The operator is starting the Pods and creating the respective Services.

## Checking peering

All nodes start with default --peering dynamic, and seeds hosts are being added for the new hosts,
to check the peers of each node you can use:

```
./scripts/peers.sh
```

The output is something like:

```
* sawtooth-1-pod-0
    tcp://service-1:8800,tcp://service-2:8800

* sawtooth-1-pod-1: 
    tcp://service-0:8800,tcp://service-2:8800

* sawtooth-1-pod-2:
    tcp://service-0:8800,tcp://service-1:8800
```

## Development

### Refresh CRD 

After changing API files you need to refresh the CRD and APIs creation with:

```
$ make refresh-crd
```

### Running locally

To run it locally you must bring a new cluster and the operator will use the default ~/.kube/config

```
$ make run 
```