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