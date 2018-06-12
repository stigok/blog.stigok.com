---
layout: post
title: "Update Kubernetes Dashboard in Azure Container Service"
date: 2018-06-12 18:00:42 +0200
categories: kubernetes azure
---

## Motivation

The Kubernetes dashboard version of the managed cluster created with Azure Container Service (ACS) is *v1.6.3*.
As of writing the latest upstream dashboard version is *1.8.3* and [Azure-compatible images][1] are being
[automatically built][2]. Read the [changelog][dashboard-releases] to see what issues has been fixed in between
the dashboard relases.

The dashboard is a *deployment* resource that is automatically installed and deployed with
an [*addon-manager*][addon-manager].

We can see the dashboard deployment resource with `kubectl`

    $ kubectl describe deployments -n kube-system kubernetes-dashboard

## Check compatible versions

Find out what your Kubernetes server version is. As of writing, the cluster created by ACS is v1.7,
as can be verified with the CLI

    $ kubectl version
    Server Version: version.Info{Major:"1", Minor:"7", GitVersion:"v1.7.7", GitCommit:"8e1552342355496b62754e61ad5f802a0f3f1fa7", GitTreeState:"clean", BuildDate:"2017-09-28T23:56:03Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"linux/amd64"}

Checking the [compatability matrix][] in the upstream dashboard repository, the latest supported
version for my cluster is v1.7.1.

## Upgrade

Check the current version of the dashboard

    $ kubectl describe pods -n kube-system kubernetes-dashboard | grep 'Image:'
        Image:          gcrio.azureedge.net/google_containers/kubernetes-dashboard-amd64:v1.6.3

Find the FQDN of the cluster master

    $ az acs list | grep fqdn

When the cluster was created, an admin username and a SSH public key was specified. Use that information
to SSH into the master in order to update the dashboard image tag.

    $ ssh adminuser@clustername.westeurope.cloudapp.azure.com
    adminuser@k8s-master-FF00FF-0:~$ sudo vim /etc/kubernetes/addons/kubernetes-dashboard-deployment.yaml

When the file has been opened for edit as described above, update the image version tag for
the *deployment* resource (terraform-like path to property: `deployment.kube-system.kubernetes-dashboard.spec.template.spec.containers.1.image`)
to `gcrio.azureedge.net/google_containers/kubernetes-dashboard-amd64:v1.7.1`.

Once the file has been saved, changes should propagate immediately. Check the pod status to verify

    $ kubectl describe pods -n kube-system kubernetes-dashboard | grep 'Image:'
        Image:          gcrio.azureedge.net/google_containers/kubernetes-dashboard-amd64:v1.7.1

## References
- https://stackoverflow.com/questions/46621566/updating-kubernetes-dashboard-image-in-a-azure-acs-k8s-cluster-is-not-getting-re

[dashboard-releases]: https://github.com/kubernetes/dashboard/releases
[addon-manager]: https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/addon-manager#addon-manager
[compatability matrix]: https://github.com/kubernetes/dashboard/wiki/Compatibility-matrix
[1]: https://console.cloud.google.com/gcr/images/google-containers/GLOBAL/kubernetes-dashboard-amd64?gcrImageListsize=50
[2]: https://console.cloud.google.com/gcr/images/google-containers/GLOBAL/kubernetes-dashboard-amd64@sha256:dc4026c1b595435ef5527ca598e1e9c4343076926d7d62b365c44831395adbd0/details/info
