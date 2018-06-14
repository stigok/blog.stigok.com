---
layout: post
title:  "Troubleshooting AWS EKS kubectl with haptio-authenticator"
date:   2018-06-13 15:41:01 +0200
categories: aws kubernetes
redirect_from:
  - /2018/06/13/aws-eks-kubectl-and-haptio-authenticator
---

## Introduction

Was setting up my first EKS cluster in AWS when I followed the [*Getting started*][eks getting started]
-guides to get it up and running.

**NOTE:** Since I followed the guide pretty closely, this post is mostly for troubleshooting, and not
a complete guide of any kind.

## Installing heptio authenticator

I installed this using `go get` package manager. As I already had Go installed it proved to be easier
than manually downloading the compiled binary. Followed the [README of the GitHub repository][heptio-repo]

    $ go get -u github.com/heptio/authenticator/cmd/heptio-authenticator-aws

And make sure that `$GOPATH/bin` is in your executable path. Output the current default `GOPATH`

    $ go env GOPATH

And verify that the path is part of your `PATH` environment variable

    $ printenv PATH

For detailed instructions, refer to the [Configure kubectl for Amazon EKS][] guide

## Troubleshooting
### `kubectl cluster-info` unauthorized

Getting an error when trying to get cluster information

```
$ kubectl cluster-info
Kubernetes master is running at https://82F83B174D696481E1CCBC3B4E817BEC.yl4.us-east-1.eks.amazonaws.com

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
error: You must be logged in to the server (Unauthorized)
```

Another hint came along when I tried to get a list available EKS clusters through the `aws` CLI

```
$ aws eks list-clusters

Could not connect to the endpoint URL: "https://eks.eu-west-1.amazonaws.com/clusters"
```

Got it solved by changing the default region to `us-east-1` (where the cluster is located).
It was initially set to `eu-west-1`, a region of which does not even support EKS yet, hence
the 404-like error message.

```
$ aws configure
AWS Access Key ID [****************x6x7]:
AWS Secret Access Key [****************xx42]:
Default region name [eu-west-1]: us-east-1
Default output format [None]
```

### heptio-authenticator-aws not found in path

Attempting to use `kubectl` with the cluster yields errors about the missing authenticator

```
$ kubectl cluster-info dump
Unable to connect to the server: getting token: exec: exec: "heptio-authenticator-aws": executable file not found in $PATH
```

### Still error unauthorized

Make sure that the command in the kubectl configuration file contains the actual cluster name as it is output by the AWS CLI

```
$ aws eks list-clusters
{
    "clusters": [
        "my-test-cluster"
    ]
}
```

Verify this name towards the command arguments passed to the `heptio-authenticator-aws`

```
$ grep -C 4 heptio ~/.kube/config
- name: my-test-cluster-admin-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: heptio-authenticator-aws
      args:
      - token
      - -i
      - my-test-cluster
```

## References
- https://stackoverflow.com/questions/50791303/kubectl-error-you-must-be-logged-in-to-the-server-unauthorized-when-accessing
- https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html
- https://github.com/heptio/authenticator#4-set-up-kubectl-to-use-heptio-authenticator-for-aws-tokens

[eks getting started]: https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html
[Configure kubectl for Amazon EKS]: https://docs.aws.amazon.com/eks/latest/userguide/configure-kubectl.html
