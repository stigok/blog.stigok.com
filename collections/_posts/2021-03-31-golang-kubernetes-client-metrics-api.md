---
layout: post
title:  "Querying metrics from the Kubernetes API using the official golang Kubernetes client"
date:   2021-03-31 18:51:41 +0200
categories: golang kubernetes
excerpt: Query the Kubernetes metrics API using the golang Kubernetes client.
#proccessors: pymd
---

## Preface

I was initially struggling a bit to figure out how to query the metrics.k8s.io API
using the golang client. For some time, it just seemed simpler to just use the client
to get a `Bearer` token and query the API directly with the `net/http`, which is
exactly what I ended up doing when I attempted the same in Python, but I wanted
to figure this one out.

## The road towards querying the Kubernetes API in golang

Within the code blocks here, I am including the required `import`s, the relevant
function's signature, and the function call itself for illustrative purposes.

### Importing auth plugins

In order to get the authentication to work properly for some cloud providers (e.g.
Azure AKS) [you have to import an auth plugin](https://github.com/kubernetes/client-go/issues/839#issuecomment-669175919).

```
import _ "k8s.io/client-go/plugin/pkg/client/auth"
```

### Configuring clients

Then we need to configure the client with a cluster (master) hostname, certificate
authority and a service account token (the actual `Bearer` token).

> Package clientcmd provides one stop shopping for building a working client from a fixed config, from a .kubeconfig file, from command line flags, or from any merged combination.
<br><small><https://pkg.go.dev/k8s.io/client-go/tools/clientcmd></small>

```go
import "k8s.io/client-go/tools/clientcmd"
// func BuildConfigFromFlags(masterUrl, kubeconfigPath string) (*restclient.Config, error)
kubeconfig, err := clientcmd.BuildConfigFromFlags("", "")
```

- The first argument here is `masterUrl`; a URL to the control plane of which API can be reached.
- The second argument is `kubeconfigPath`; the location of a Kubernetes configuration file, just like
  the ones you normally have locally on your computer when using `kubectl`.

If both these arguments are empty strings, it will fall back to *in-cluster config*,
which will configure it to use an in-cluster reachable `masterUrl` and look for
CA and token in the default locations (`/run/secrets/kubernetes.io/serviceaccount/{ca.crt,token}`).

So when you are running stuff off of your local machine, which has a working kubeconfig,
pass an empty `masterUrl` and use the location of your config file to `kubeconfigPath`,
which is normally at `~/.kube/config`.

### Instantiating an API client

Now we need to instantiate a client set using the configuration we have loaded.
A client set contains API groups, and the client set exported in
[`k8s.io/client-go/kubernetes`][client-go-kubernetes] contains all the officially
supported groups for the chosen package version.
This means that the latest version, as of writing, includes all the API groups
available in Kubernetes v1.20.

For the metrics API however, we need to use an additional client set which is defined in
[`k8s.io/metrics/pkg/client/clientset/versioned`][k8s-metrics].

```go
import "k8s.io/client-go/kubernetes"
// func NewForConfigOrDie(c *rest.Config) *Clientset
clientset := kubernetes.NewForConfigOrDie(kubeconfig)

import metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
// func NewForConfigOrDie(c *rest.Config) *Clientset
metricsclientset := metricsv.NewForConfigOrDie(kubeconfig)
```

The `NewForConfigOrDie` is a helper function to `panic` if an error occurs.
The `NewForConfig` can be used instead if you want to handle the `error` manually.

### Querying using the client sets
#### The CoreV1 API

Now we are ready to send actual queries. As noted, the clientsets contains one or more
API groups. To query pods, we have to use the `CoreV1` group. To figure out what group you
have to use for other resources I advice to read the [Kubernetes API reference][api-ref].

Now, having a `kubeconfig` and a `clientset` we can list all pods in all namespaces.
The single string argument passed to `Pods()` is a namespace. When empty it matches
all namespaces.

```go
pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
```

- The `context.TODO()` is a placeholder for a context. It allows for
  cancelling a pending requests if used with a `context.WithCancel`.
- The [`metav1.ListOptionts{}`][metav1-list-options] allows for query parameters like label selectors,
  a `Watch` property to get an event stream, specifying desired API resource versions returned
  etc.

What this query returns is a [`*v1.PodList`][v1-pod-list]. To see if any pods were found
you can check the `len` of `PodList.Items`, or iterate over it with `range`.

See [the godoc for v1.Pod](https://pkg.go.dev/k8s.io/api/core/v1#Pod) to see what
properties this struct exports.
The struct is importing other structs inline so e.g. the metadata of a resource is at the
root level of the struct itself, not under a `Metadata` property. I.e. the name and
namespace of a `Pod` is in `Pod.Name` and `Pod.Namespace`.

#### Metrics API

For this we use the other client set.

```go
podMetrics, err := metricsclientset.MetricsV1beta1().PodMetricses("").List(context.TODO(), metav1.ListOptions{})
```

See the godoc for the [metrics v1beta1][] for details of how the resource is
structured. It's pretty much the same as with the `Pod`s, except it's a different
resource with different properties.

## Afterwords

I think this post sheds light on the complexity of the Kubernetes project.
Having to use multiple packages and modules
to do a seemingly simple operations with the Kubernetes API is drastically raising the
bar of entry. I can't say I know how things should be done differently,
but here's an attempt to help someone else out there to do the same thing as I was
striving for; getting all pods and pod metrics in all namespaces in a Kubernetes
cluster.

## References

[clientcmd]: https://pkg.go.dev/k8s.io/client-go/tools/clientcmd
[client-go-kubernetes]: https://pkg.go.dev/k8s.io/client-go/kubernetes
[k8s-metrics]: https://pkg.go.dev/k8s.io/metrics/pkg/client/clientset/versioned
[metav1-list-options]: https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#ListOptions
[metrics v1beta1]: https://pkg.go.dev/k8s.io/metrics@v0.20.5/pkg/client/clientset/versioned/typed/metrics/v1beta1#PodMetricsInterface
[api-ref]: https://v1-19.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/
[v1-pod-list]: https://pkg.go.dev/k8s.io/api/core/v1#PodList
