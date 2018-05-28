# tiller-releases-convertor

It is a tool to automate [Helm](http://helm.sh/)'s [tiller](https://docs.helm.sh/glossary/#tiller) migration from ConfigMap releases backend to Secrets-based backend for Kubernetes & Helm users.

## TL;DR

To upgrade your Helm's Tiller setup to Secrets storage backend with zero-downtime run these three commands in the following order:

```shell
tiller-releases-convertor convert # Will create a Secret for every ConfigMap release
tiller-releases-convertor secure-tiller # Will update tiller deployment with "--storate=secret"
tiller-releases-convertor cleanup # Deletes tiller's ConfigMaps
```

## Overview

If you use [Kubernetes](https://kubernetes.io/) container orchestrator, there is a high probability you are already using [Helm](http://helm.sh/) - the most popular "package manager for Kubernetes", which consists of two components: helm client and tiller.

Tiller is the in-cluster Helm's component, which manages Helm releases creating, updating or removing Kubernetes resources. It also manages Helm releases data. Tiller can use one of three storage types to store releases: ConfigMaps, Secrets, and Memory. The default option is ConfigMaps. If you have a default Tiller setup(with some software installed via Helm) on your cluster you can see them for yourself by typing:

```shell
kubectl get configmaps -n kube-system -l OWNER=TILLER
```

You will see all ConfigMaps, which are owned by Tiller.

When you install software using Helm, you provide it with additional Helm Chart Values to configure your installation. These values can contain a lot of sensitive data you wish to be stored as Kubernetes Secrets. ConfigMaps are stored in the Kubernetes' Etcd key-value storage as is, without any encryption. Secrets, on the other hand, are typically encrypted. The resulting Kubernetes configuration, which Helm produces, may contain any Kubernetes resources, including Secrets. But, in any case, Helm creates an in-cluster object for each version of each release, and that object is stored in a ConfigMap by default.

It's a security flaw, and there is an option to use Secrets as a storage backend for Helm.

If you make a new Tiller installation you can use the following command to enable Secrets storage backend:

```shell
helm init --service-account tiller --override 'spec.template.spec.containers[0].command'='{/tiller,--storage=secret}'
```

But what to do if you already have a deafult working installation? Basically, you have one option to export all Tiller's ConfigMaps, convert them to appropriate Secrets, apply the latter and edit Tiller Deployment by hands. This can be a challenge if you have a lot of releases in your cluster.

This one-time script was written to solve that problem for you.

This tool has four commands:

* [list](#list)
* [convert](#convert)
* [secure-tiller](#secure-tiller)
* [cleanup](#cleanup)

The normal way to use it looks like that:

```shell
tiller-releases-convertor convert # Will create a Secret for every ConfigMap release
tiller-releases-convertor secure-tiller # Will update tiller deployment with "--storate=secret"
tiller-releases-convertor cleanup # Deletes tiller's ConfigMaps
```

## Supported Kubernetes versions

[kubernetes/client-go version 7.0.0](https://github.com/kubernetes/client-go) is used in this script. See original [Compatibility matrix](https://github.com/kubernetes/client-go#compatibility-matrix).

## Installation

It is strongly recommended that you use a released version. You can find current release binaries on the [releases](https://github.com/dragonsmith/tiller-releases-converter/releases) page.

Or you can install a classical way:

```shell
go get -u https://github.com/dragonsmith/tiller-releases-converter
```

## Commands overview

### list

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

### convert

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

See "target already exists" error? It happens when the target Secret already exists in the cluster. This tool will just skip this Secret, and you shall check every situation like that by yourself.

### secure-tiller

This command updates `tiller-deploy` deployment and adds `.spec.template.spec.containers[0].command={"/tiller", "--storage=secret"}`

```shell
tiller-releases-convertor secure-tiller
```

Output example:

```
Updating Tiller Deployment...
Tiller Deployment was updated successfully!
```

### cleanup

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

Red cross indicates an error which is followed by an error message.

## Command line options

There are some global command line arguments:

```
Flags:
      --context string      kube config context
  -h, --help                help for tiller-releases-convertor
  -c, --kubeconfig string   config file (default is $HOME/.kube/config)
  -n, --namespace string    tiller namespace (default is kube-system)
```

## Contributing

1. Fork it
2. Download your fork  (`git clone https://github.com/your_username/tiller-releases-converter && cd tiller-releases-converter`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Make changes and add them (`git add .`)
5. Commit your changes (`git commit -m 'Add some feature'`)
6. Push to the branch (`git push origin my-new-feature`)
7. Create a new pull request

## License
tiller-releases-convertor is released under the Apache 2.0 license. See [LICENSE](https://github.com/dragonsmith/tiller-releases-converter/blob/master/LICENSE)
