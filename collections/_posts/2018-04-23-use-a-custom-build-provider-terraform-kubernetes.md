---
layout: post
title:  "Use a custom built provider with Terraform"
date:   2018-04-22 19:48:07 +0200
categories: terraform kubernetes
---

I was using a custom fork of the `kubernetes` provider for terraform because
the upstream version does not have support for `ReplicaSets` or `StatefulSets`,
(and some more) which is something I needed for a Kafka deployment.

Terraform is built with Go, so are the providers.
I did not have `go` installed before I started, so I installed that first

    $ yaourt -S go

Then creating the default Go path directory

    $ mkdir $(go env GOPATH)

Set environment for the next couple of commands

    $ export GOPATH="$(go env GOPATH)"

Cloned the fork into this folder as described in the [fork's README.md](fork-readme),
and build from source

    $ mkdir -p $GOPATH/src/github.com/sl1pm4t; cd $GOPATH/src/github.com/sl1pm4t
    $ git clone git@github.com:sl1pm4t/terraform-provider-kubernetes
    $ make build

A binary should now have been built and placed in `$GOPATH/bin`. Use this path when initialising
the terraform environment in your project

    $ cd ~/projects/tf-test
    $ terraform init -plugin-dir=$GOPATH/bin

Assuming I have a `.tf` file in that folder containing a `provider` resource looking something like this

```terraform
provider "kubernetes" {}
```

It should now be using the locally built kubernetes provider instead of the upstream one.

[fork-readme]: https://github.com/sl1pm4t/terraform-provider-kubernetes/blob/custom/README.md
