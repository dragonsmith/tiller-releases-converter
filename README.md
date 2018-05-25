# tiller-releases-convertor

It is a tool to automate [Helm](http://helm.sh/)'s [tiller](https://docs.helm.sh/glossary/#tiller) migration from ConfigMap releases backend to Secrets-based backend.

## TL;DR

To upgrade your Helm's Tiller setup to Secrets storage backend with zero-downtime run these three commands in the following order:

```shell
tiller-releases-convertor convert # Will create a Secret for every ConfigMap release
tiller-releases-convertor secure-tiller # WIll update tiller deployment with "--storate=secret"
tiller-releases-convertor cleanup # Deletes tiller's ConfigMaps
```

## Overview

Tiller is the in-cluster Helm's component, which manages Helm releases creating, updating or removing Kubernetes resources. It also manages Helm releases data. Tiller can use one of three storage types to store releases: ConfigMaps, Secrets, and Memory. The default option is ConfigMaps. If you have a default Tiller setup(with some software installed via Helm) on your cluster you can see them for yourself by typing:

```shell
kubectl get configmaps -n kube-system -l OWNER=TILLER
```

You will see all ConfigMaps, which are owned by Tiller.

Here comes the main problem. These releases contain every value you've passed during helm install/update including your app's secrets. ConfigMaps are stored in the Kubernetes etcd key-value storage as is, without any encryption. Secrets, on the other hand, are normally encrypted.

If you make a new Tiller installation it's recommended use the next command to enable Secrets storage backend:

```shell
helm init --service-account tiller  --override 'spec.template.spec.containers[0].command'='{/tiller,--storage=secret}'
```
But what to do if you have a working installation?

This script solves the problem. It


This tool has four commands:

* list
* convert
* secure-tiller
* cleanup

## list

Lists current tiller's ConfigMap-based releases.

```shell
tiller-releases-convertor list
```

Output example:
```
I've found these Tiller's ConfigMap releases for you:

kube-state-metrics.v1
kube-state-metrics.v2
...
```

## convert

This command creates Secrets out of Tiller's ConfigMaps. ConfigMaps are left untouched.

```shell
tiller-releases-convertor convert
```

Output Example:

```
 - [🚫] kube-state-metrics.v1 (target already exists)
 - [✅] kube-state-metrics.v2
...
```

See "target already exists" error? It happens when the target Secret already exists in the cluster. This tool will just skip this Secret, you shall check every situation like that by yourself.

## secure-tiller

This command updates `tiller-deploy` deployment and adds `.spec.template.spec.containers[0].command={"/tiller", "--storage=secret"}`

```shell
tiller-releases-convertor secure-tiller
```

Output example:

```
Updating Tiller Deployment...
Tiller Deployment was updated successfully!
```

## cleanup

This command deletes old Tiller ConfigMaps.

```shell
tiller-releases-convertor cleanup
```

Output example:

```
Deleting: kube-state-metrics.v1 ✅
Deleting: kube-state-metrics.v2 ❌
<Error message>

```

Red cross indicates an error which is followed by en error message.

## Command line options

There are some global command line arguments:

```
Flags:
      --context string      kube config context
  -h, --help                help for tiller-releases-convertor
  -c, --kubeconfig string   config file (default is $HOME/.kube/config)
  -n, --namespace string    tiller namespace (default is kube-system)
```
