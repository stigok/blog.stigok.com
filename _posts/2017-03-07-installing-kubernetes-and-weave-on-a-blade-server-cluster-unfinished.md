---
layout: post
title: "Installing kubernetes and weave on a blade server cluster (unfinished)"
date: 2017-03-07 19:35:50 +0100
categories: kubernetes linux draft
redirect_from:
  - /post/installing-kubernetes-and-weave-on-a-blade-server-cluster-unfinished
---

**Goal:** Get 28 blade servers up and running in a kubernetes cluster

# NOTE

> This article is not finished, and will change.

# Planning

Things (I think) I have to do

- Set up [Weave Net](https://www.weave.works/docs/net/latest/introducing-weave/) that _creates a virtual network that connects Docker containers across multiple hosts and enables their automatic discovery_
- Install Kubernetes on one of the blades
- Install Kubernetes on an additional blade and add it to the cluster
- Verify the child blade configuration and create a pre-configured ISO that can be easily deployed to the rest of the blades
- Figure out what to use the cluster for

## Setting up Weave Net

The reason this is the first step, before even touching the Kubernetes installation, is
 - we don't have any spare public IPv4 addresses for the blades
 - we have an ocean of free IPv6 addresses, but Kubernetes does not play well with IPv6-only hosts
 - when different people are configuring their own containers, we want the networks to be isolated
 - we want to centralize the network configuration without involving additional hardware

#### Installing Weave

This is the machine I'm starting off with

    # cat /etc/issue.net 
    Ubuntu 16.04.2 LTS

Verify that the version of docker is > 1.6

    # docker --version
    Docker version 17.03.0-ce, build 3a232c8

Install weave 

    # curl -L git.io/weave -o /usr/local/bin/weave
    # chmod a+x /usr/local/bin/weave

I want to disable the automatic version check

    # export CHECKPOINT_DISABLE=1

Time to sit down and read for a bit: [Understanding how Weave Net Works](https://www.weave.works/docs/net/latest/how-it-works/)

Launch weave on the main host

    # weave launch

After launch, it should have pulled all of its docker images

    # weave version
    weave script 1.9.3
    weave router 1.9.3
    weave proxy  1.9.3
    weave plugin 1.9.3

Ensure that containers launched via the Docker command line are automatically attached to the Weave network

    # eval $(weave env)

Start a container and make sure no errors are being thrown upon creation. (Read [Using Weave Net](https://www.weave.works/docs/net/latest/using-weave/) for troubleshooting errors)

    # docker run --name test-1 -it weaveworks/ubuntu

Success!

#### Blade server #1

This machine is only reachable through IPv6 and therefore needs some extra tender love and care. Add an IPv6 DNS server address (Google DNS)

    # echo "nameserver 2001:4860:4860:0:0:0:0:8888" >> /etc/resolv.conf

Then see if it works

    # dig +short stigok.com AAAA
    2a02:ed06::211

Now another test

    # dig +short git.io AAAA

Did not give us anything. Which means I have to find another way to download the `weave` binary on this host. However related, [Where is the git.io link redirecting to?], it did not solve this issue, as github.com itself doesn't have any AAAA records either. I'm using SCP to transfer the binary directly. Note that it has to be put in a directory in your `$PATH`. Since this box is Fedora, `/usr/local/bin` is not by default. However, `/usr/local/sbin` is.

    # scp root@weave-host:/usr/local/bin/weave /usr/local/sbin/weave

Make sure docker is running

    # systemctl status docker

If not, start it

    # systemctl start docker

The docker image repository is not IPv6 enabled, so grab the images from the other host. To enable my non-root user to control the Docker daemon, I add myself to the `docker` group. On Fedora, this group is not added by default, so first I create it.

    # groupadd docker
    # gpasswd --add user docker

List images 

{% raw %}
    # ssh user@weave-host docker images --format '{{.Repository}}'
    weaveworks/plugin
    weaveworks/weave
    weaveworks/weaveexec
    weaveworks/weavedb
    weaveworks/ubuntu
{% endraw %}

Ok, so now lets fetch all of those and load them into this host

    # ssh user@weave-host docker save weaveworks/plugin weaveworks/weave weaveworks/weaveexec weaveworks/weavedb weaveworks/ubuntu | docker load
    7cbcbac42c44: Loading layer [==================================================>]  5.05 MB/5.05 MB
    43f5e1d45225: Loading layer [==================================================>] 2.048 kB/2.048 kB
    cdd5598a1676: Loading layer [==================================================>] 18.97 MB/18.97 MB
    b7828b5e45cc: Loading layer [==================================================>] 12.94 MB/12.94 MB
    7f2483959486: Loading layer [==================================================>] 9.478 MB/9.478 MB
    5ca8af2897f9: Loading layer [==================================================>] 4.581 MB/4.581 MB
    a78e0bf11eb8: Loading layer [==================================================>] 3.206 MB/3.206 MB
    782b8b79808c: Loading layer [==================================================>] 4.512 MB/4.512 MB
    fbbaf58e6c6a: Loading layer [==================================================>]  34.9 MB/34.9 MB
    86d0e9420b70: Loading layer [==================================================>] 26.48 MB/26.48 MB
    ead8a0e246bb: Loading layer [==================================================>] 2.048 kB/2.048 kB
    fe4aaeadcb92: Loading layer [==================================================>]   277 kB/277 kB
    45698d068987: Loading layer [==================================================>] 9.681 MB/9.681 MB
    Loaded image: weaveworks/plugin:1.9.3
    Loaded image: weaveworks/weave:1.9.3
    Loaded image: weaveworks/weaveexec:1.9.3
    d8330d2726bf: Loading layer [==================================================>] 2.048 kB/2.048 kB
    Loaded image: weaveworks/weavedb:latest
    9436069b92a3: Loading layer [==================================================>] 127.6 MB/127.6 MB
    19429b698a22: Loading layer [==================================================>] 14.85 kB/14.85 kB
    82b57dbc5385: Loading layer [==================================================>] 11.78 kB/11.78 kB
    737f40e80b7f: Loading layer [==================================================>] 4.608 kB/4.608 kB
    5f70bf18a086: Loading layer [==================================================>] 1.024 kB/1.024 kB
    a7cf6da36ea3: Loading layer [==================================================>] 3.968 MB/3.968 MB
    Loaded image: weaveworks/ubuntu:latest

Verify that the images are now locally available

{% raw %}
    # docker images --format '{{.Repository}}'
    weaveworks/plugin
    weaveworks/weave
    weaveworks/weaveexec
    weaveworks/weavedb
    weaveworks/ubuntu
{% endraw %}

Launch weave and connect it to the weave host

    # weave launch weave-host

Allow communication on the non-weave network, to avoid blocking communication between the weave nodes

    iptables -I INPUT -s 10.8.1.0/24 -j ACCEPT

### This article is not finished. Come back later
