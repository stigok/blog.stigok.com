---
layout: post
title: "kubectl proxy error connection refused and too many open files"
date: 2018-02-27 00:28:22 +0100
categories: gce kubernetes
redirect_from:
  - /post/kubectl-proxy-error-connection-refused-and-too-many-open-files
---

I just got into testing out a managed kubernetes (k8s) cluster on Google Cloud Engine and got into some early trouble.

    $ kubectl proxy
    Starting to serve on 127.0.0.1:8001
    I0226 23:59:15.574452    2842 logs.go:41] http: proxy error: dial tcp [::1]:8080: getsockopt: connection refused

Then, that error message made me thing I should maybe be listening on port 8080 instead

    $ kubectl proxy 
    Starting to serve on 127.0.0.1:8080
    I0226 23:00:24.463088    2885 logs.go:41] http: Accept error: accept tcp 127.0.0.1:8080: accept4: too many open files; retrying in 5ms
    I0226 23:00:24.463677    2885 logs.go:41] http: proxy error: dial tcp 127.0.0.1:8080: socket: too many open files

Apparently, this is what happens when kubectl doesn't have a valid config yet. 

    $ kubectl config view
    apiVersion: v1
    clusters: []
    contexts: []
    current-context: ""
    kind: Config
    preferences: {}
    users: []

It doesn't matter if the `gcloud` instance has a good connection - `kubectl` needs its own.

If `gcloud` has not yet been set up, do that first (`gcloud init`), then load credentials for `kubectl`:

    $ gcloud container clusters list
    NAME       LOCATION        MASTER_VERSION  MASTER_IP       MACHINE_TYPE  NODE_VERSION  NUM_NODES  STATUS
    cluster-1  europe-west1-c  1.8.7-gke.1     185.0.0.42  g1-small      1.8.7-gke.1   4          RUNNING

And I want to operate on `cluster-1`, so I can get credentials for that one by name:

    $ gcloud container clusters get-credentials cluster-1
    Fetching cluster endpoint and auth data.
    kubeconfig entry generated for cluster-1.

Success! Now I can start the proxy

    $ kubectl proxy
    Starting to serve on 127.0.0.1:8001

And browse the web interface of the cluster at `http://127.0.0.1:8001/ui`