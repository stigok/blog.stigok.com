---
layout: post
title:  "Terraform provider kubernetes - system:anonymous cannot create services in the namespace"
date:   2018-06-25 18:21:05 +0200
categories: aws eks kubernetes terraform
---

## Introduction
Had problems trying to apply kubernetes resources using the [terraform
provider for kubernetes][] on a AWS EKS Kubernetes cluster. I was getting
errors I had not previously had before with neither Azure AKS nor Google GKE
clusters:

```
Error: Error applying plan:

15 error(s) occurred:

* kubernetes_service.pzoo: 1 error(s) occurred:

* kubernetes_service.pzoo: services is forbidden: User "system:anonymous" cannot create services in the namespace "test-kafka"
* kubernetes_service.zoo: 1 error(s) occurred:

* kubernetes_service.zoo: services is forbidden: User "system:anonymous" cannot create services in the namespace "test-kafka"
* kubernetes_secret.container-registry: 1 error(s) occurred:
*

[redacted]
```

This apparently has something to do with the EKS cluster is using roles based
authentication (RBAC), and there was not a proper service account set in my
current `kubectl` context.

## Solution

It is expected that there is already a cluster definition in the current
kubectl configuration. Check out my other post about [getting authenticated
using heptio-authenticator](https://blog.stigok.com/2018/06/13/aws-eks-kubectl-and-heptio-authenticator.html)
for more information.

Create a new service account for terraform

    $ kubectl -n kube-system create sa terraform
    $ kubectl create clusterrolebinding terraform --clusterrole cluster-admin
    --serviceaccount=kube-system:terraform

Get the bearer authorization token for the created user and save the credentials
in the kubectl configuration

    $ kubectl config set-credentials terraform --token "$(kubectl describe secrets -n kube-system terraform | grep -Po '^token:\s+\K\S+$')"

Create a new context for the EKS cluster with the above credentials for terraform to use

    $ kubectl config set-context aws-test-terraform --cluster=my-aws-cluster --user=terraform

Verify that it works as it should (the below command should indeed output *yes*)

    $ kubectl auth can-i --context=aws-test-terraform create deployments
    yes

Specify the context in the terraform kubernetes provider definition

```
provider "kubernetes" {
  config_context = "aws-test-terraform"
}
```

Now use `terraform apply` as normal.

## References

- https://dzone.com/articles/terraform-vs-helm-for-kubernetes

[terraform provider for kubernetes]: https://github.com/sl1pm4t/terraform-provider-kubernetes
