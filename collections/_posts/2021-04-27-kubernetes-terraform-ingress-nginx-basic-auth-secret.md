---
layout: post
title:  "Protecting a Kubernetes ingress for ingress-nginx with HTTP basic auth using Terraform"
date:   2021-04-27 16:54:22 +0200
categories: kubernetes nginx terraform
excerpt: Apply basic auth to a Kubernetes Ingress for ingress-nginx using Terraform
#proccessors: pymd
---

## Preface

I am working in a multi-cluster environment now where we are unable to see the real client IP.
The Kubernetes clusters are behind "dumb" load balancers for ingress traffic which forwards traffic
to our ingress-nginx service endpoints. When traffic arrives, the source IP stems from the LB,
and the `X-Forwarded-For` header cannot be trusted. This means we cannot use ingress-nginx's
`nginx.ingress.kubernetes.io/whitelist-source-range` to protect our `Ingress` resources from
the public. Instead we will rely on authentication.

**Please note the caveats near the end of this post.**

## Setup

If you're using ingress-nginx in Kubernetes, you can [apply HTTP basic auth using annotations](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/#authentication).

For this example I am using the `auth-map` as the `auth-secret-type`, which reads a `Secret`'s `data` fields
and uses its keys as usernames and its values as the hashed passwords.

```yaml
nginx.ingress.kubernetes.io/auth-type: basic
nginx.ingress.kubernetes.io/auth-secret: my-namespace/my-auth-secret
nginx.ingress.kubernetes.io/auth-secret-type: auth-map
```

Now, the secret in `my-namespace` can be populated with `username: <base64 encoded bcrypt hash>`.
Since this is also about Terraform, I'll be creating the secret using Terraform.

```terraform
locals {
  app_namespace = "my-namespace"
  app_name      = "my-app"
  auth_username = "stigok"
}

resource "random_password" "ingress-auth" {
  length           = 32
  special          = false
  override_special = ",.-_!"
}

resource "kubernetes_secret" "ingress-auth" {
  metadata {
    name      = "${local.app_name}-basic-auth"
    namespace = local.namespace
    labels = {
      app       = local.app_name
    }
  }
  data = {
    local.auth_username = bcrypt(random_password.ingress-auth.result)
  }
}

resource "kubernetes_ingress" "my-app" {
  metadata {
    name      = local.app_name
    namespace = local.app_namespace
    annotations = {
      "nginx.ingress.kubernetes.io/auth-type"        = "basic",
      "nginx.ingress.kubernetes.io/auth-secret-type" = "auth-map",
      "nginx.ingress.kubernetes.io/auth-secret"      = "${local.app_namespace}/${kubernetes_secret.ingress-auth.metadata.0.name}"
      "nginx.ingress.kubernetes.io/auth-realm"       = "auth required for ${local.app_name}"
    }
    labels = {
      app = local.app_name
    }
  }

  spec {
    rule {
      host = local.app_hostname
      http {
        path {
          path = "/"
          backend {
            service_name = "my-service"
            service_port = "8080"
          }
        }
      }
    }

    tls {
      hosts = [local.app_hostname]
    }
  }
}
```

## Caveats

However, this has some implications...

### cert-manager
This requires auth for all paths for the host which will make cert-manager unable to validating certificate requests
using HTTP01 validation (paths under `/.well-known`).
So if you rely on ACME HTTP01 for aquiring certificates for a single host then this will break your setup.
You should be good if you have an existing wildcard certificate installed in your cluster or if you're
using DNS01 validation.

### External monitoring solutions
It also blocks external `/health` requests. So if you have a monitoring solution outside your cluster that
needs unauthenticated access to specific routes, you can create an additional ingress without auth, that
only matches the `/health` endpoint. If using the `nginx.ingress.kubernetes.io/use-regex: "true"`, you can
make your ingress match a single route only using `/health$` as the path in the ingress spec.

However, if your monitoring solution allows you to provide basic auth credentials this will not be a problem.
You can create new set of credentials for it and have it authenticate like all other clients.

### State is always dirty
A big pain here is that `bcrypt` uses a randomly selected salt value causing it to return a new hash on
every call. This has the effect of making your state dirty on each and every plan. For those who strives
to always have a clean state (like me) this is a tad annoying. I haven't figured out a work-around yet.

## References
- <https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/annotations/>
- <https://www.terraform.io/docs/language/functions/bcrypt.html>
