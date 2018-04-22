---
layout: post
title: "Terraform secret creation invalid character looking for beginning of value"
date: 2018-04-09 02:11:51 +0200
categories: terraform kubernetes docker
redirect_from:
  - /post/terraform-secret-creation-invalid-character-looking-for-beginning-of-value
---

I was trying to create a kubernetes_secret with Terraform when I got this error

    * kubernetes_secret.container-repository: Secret "azurecr" is invalid: data[.dockerconfigjson]: Invalid value: "<secret contents redacted>": invalid character 'e' looking for beginning of value

I read in the [Kubernetes container image reference](https://kubernetes.io/docs/concepts/containers/images/#bypassing-kubectl-create-secrets) that I would set the value of the secrets as the base64 representation of a Docker `config.json` file:

    apiVersion: v1
    kind: Secret
    metadata:
      name: myregistrykey
      namespace: awesomeapps
    data:
      .dockerconfigjson: UmVhbGx5IHJlYWxseSByZWVlZWVlZWVlZWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGxsbGx5eXl5eXl5eXl5eXl5eXl5eXl5eSBsbGxsbGxsbGxsbGxsbG9vb29vb29vb29vb29vb29vb29vb29vb29vb25ubm5ubm5ubm5ubm5ubm5ubm5ubm5ubmdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2cgYXV0aCBrZXlzCg==
    type: kubernetes.io/dockerconfigjson

However, this is not the case when declaring them in terraform with `kubernetes_secret`, in which case you would avoid encoding it. Instead, you can include it with the `file()` interpolation syntax

    resource "kubernetes_secret" "container-repository" {
      metadata {
        name = "azurecr"
        namespace = "${var.namespace}"
      }
    
      data {
        ".dockerconfigjson" = "${file("${path.module}/.docker/config.json")}"
      }
    
      type = "kubernetes.io/dockerconfigjson"
    }


## References
- https://kubernetes.io/docs/concepts/containers/images/#bypassing-kubectl-create-secrets
- https://www.terraform.io/docs/providers/kubernetes/r/secret.html
- https://github.com/terraform-providers/terraform-provider-kubernetes/issues/145