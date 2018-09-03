---
layout: post
title:  "Define imagePullSecrets on the default service account in namespace"
date:   2018-09-03 17:43:05 +0200
categories: kubernetes
---

## Introduction

I want to avoid annotating all my pods and pod template manifests with `imagePullSecrets`.

When `imagePullSecrets` hasn't been set, the secrets of the default service account in the current namespace is used instead.
If those aren't defined either, default or no credentials are used, unless something provider specific is happening.

AFAIK, in AWS, the default service account has access to the AWS ECR based through AWS IAM.

## Configuring

Create a secret containing the container registry credentials **in the same namespace** as they will be used.

```
$ export NAMESPACE=myproject SERVER=myprojectcr.azurecr.io USERNAME=user PASSWORD='0226+1111' NAME=myproject-container-registry
$ kubectl create secret -n $NAMESPACE docker-registry --docker-server=$SERVER --docker-username=$USERNAME --docker-password="$PASSOWRD" $NAME
```

Patch the default service account in **the same namespace** with the `imagePullSecrets`

```
$ kubectl patch sa default -n $NAMESPACE -p '"imagePullSecrets": [{"name": "myproject-container-registry" }]'
```

The service account should now be able to authenticate itself with the remote container registry and we no longer have to specify them on each `Pod` or `Template` definition.

## References
- https://stackoverflow.com/questions/40288077/how-to-pass-image-pull-secret-while-using-kubectl-run-command

