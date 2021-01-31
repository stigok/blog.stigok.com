---
layout: post
title:  "Configure Kubernetes deployment secrets in GitHub Actions with Terraform"
date:   2021-01-31 15:25:20 +0100
categories: terraform kubernetes
excerpt: Configuring Kubernetes deployment secrets in GitHub Actions with Terraform
#proccessors: pymd
---

## Preface

I've recently been using GitHub actions a lot and often deploying to Kubernetes in the GitHub Actions workflows.
The Kubernetes workflow step uses [`Azure/k8s-set-context`](https://github.com/Azure/k8s-set-context/)
to connect and authenticate with the cluster.

```
- uses: azure/k8s-set-context@v1
  with:
    method: service-account
    k8s-url: <URL of the clsuter's API server>
    k8s-secret: <secret associated with the service account>
  id: setcontext
```

This expects a Kubernetes service account configured in the cluster, preferrably in
a single target namespace. This requires a `Role`, `RoleBinding` and a `ServiceAccount`. It's fairly
simple to set up manually when you have a single repo, but when you have to do this manually on a regular
basis, it gets tedious and error prone.

The Terraform configuration I'm presenting here will set up all the required Kubernetes resources
in the desired namespace, and put the resulting `ServiceAccount` secret token as a secret in the GitHub
repository.

Although this configuration works, it is simplified for legibility.
Make sure you are protecting your secrets by turning them into `sensitive` `variable`s.

## Terraform Configuration

These are the versions I've tested with. The Kubernetes is now in 2.x, but I haven't upgraded yet (might work out of the box).

```terraform
terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = "~> 1.13"
    }
    github = {
      source  = "integrations/github"
      version = "~> 4.3"
    }
  }
  required_version = ">= 0.14"
}
```

Please refer to the provider documentation if you're not familiar with configuring them.
- <https://registry.terraform.io/providers/hashicorp/kubernetes/1.13.3/docs>
  - You can either use your local `$KUBECONFIG` or get credentials from your cloud provider's own Terraform provider.
- <https://registry.terraform.io/providers/integrations/github/latest/docs>
  - You can get a token from <https://github.com/settings/tokens>

```terraform
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 1.13"
    }
    github = {
      source  = "integrations/github"
      version = "~> 4.3"
    }
  }
  required_version = ">= 0.14"
}

locals {
  name = "github-actions"

  github_repo_owner = "my-org-or-username"
  github_repo_name  = "my-repo"

  kubernetes_namespace = "my-namespace"
  kubernetes_api_url   = "can be found in your kubeconfig or from your cloud provider"
}

provider "kubernetes" {}

provider "github" {
  owner = local.github_repo_owner
  token = "a valid personal access token"
}

resource "kubernetes_role" "actions" {
  metadata {
    name      = local.name
    namespace = local.kubernetes_namespace
  }

  rule {
    api_groups = ["*"]
    resources  = ["*"]
    verbs      = ["get", "list", "watch", "create", "update", "patch"]
  }
}

resource "kubernetes_service_account" "actions" {
  metadata {
    name      = local.name
    namespace = local.kubernetes_namespace
  }

  depends_on = [
    kubernetes_role.actions,
  ]
}

data "kubernetes_secret" "actions" {
  metadata {
    name      = kubernetes_service_account.actions.default_secret_name
    namespace = local.kubernetes_namespace
  }

  depends_on = [
    kubernetes_service_account.actions,
  ]
}

resource "kubernetes_role_binding" "actions" {
  metadata {
    name      = local.name
    namespace = local.kubernetes_namespace
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = local.name
  }
  subject {
    kind      = "ServiceAccount"
    name      = local.name
    namespace = local.kubernetes_namespace
  }

  depends_on = [
    kubernetes_role.actions,
    kubernetes_service_account.actions,
  ]
}

resource "github_actions_secret" "kubernetes-sa" {
  repository  = local.github_repo_name
  secret_name = "K8S_SA_SECRET"
  # Mimic a Kubernetes secret in YAML.
  # The GitHub action k8s-set-context only reads the `data` field anyway.
  plaintext_value = yamlencode({
    data = {
      "ca.crt"  = base64encode(data.kubernetes_secret.actions.data["ca.crt"])
      token     = base64encode(data.kubernetes_secret.actions.data.token)
      namespace = base64encode(data.kubernetes_secret.actions.data.namespace)
    }
  })
}

resource "github_actions_secret" "kubernetes-api-url" {
  repository      = local.github_repo_name
  secret_name     = "K8S_API_URL"
  plaintext_value = local.kubernetes_api_url
}
```

To test for yourself, put all the above Terraform snippets into a single file and run
- `terraform init`
- `terraform validate`
- `terraform plan`

This is not a lot of code, especially not if you have a lot of your infrastructure configured in Terraform already.
If you make this into a module and set `variable` definitions instead of `locals`, the reusability is great!

In our GitHub Actions workflow file we can now connect to the Kubernetes cluster with the following workflow step

{% raw %}
```yaml
- uses: azure/k8s-set-context@v1
  with:
    method: service-account
    k8s-url: ${{ secrets.K8S_API_URL }}
    k8s-secret: ${{ secrets.K8S_SA_SECRET }}
  id: setcontext
```
{% endraw %}

## References
- <https://github.com/Azure/k8s-set-context/>
- <https://www.terraform.io/docs/language/>
- <https://registry.terraform.io/providers/integrations/github/latest/docs>
- <https://registry.terraform.io/providers/hashicorp/kubernetes/1.13.3/docs>
