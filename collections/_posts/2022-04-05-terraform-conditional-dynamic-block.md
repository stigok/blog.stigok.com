---
layout: post
title:  "Using conditional dynamic blocks in Terraform"
date:   2022-04-05 12:27:25 +0200
categories: terraform
excerpt: Dynamic blocks do not have count, but a hacked for_each will do.
#proccessors: pymd
---

## Preface

I wanted to use a conditional [`dynamic` block](https://www.terraform.io/language/expressions/dynamic-blocks) in my Terraform configuration,
but `dynamic` does not support `count`.

## Conditional dynamic block

Instead of `count` we can use `for_each` with a conditional map using `merge`:

```terraform
dynamic "env" {
  # A bogus map for a conditional block
  for_each = merge(var.enable_vault ? {} : { vault_disabled = true })

  content {
    name  = "DOCKER_CONFIG"
    value = "/dockerconfig"
  }
}
```

For a full example, here's an example configuration for a [Kubernetes pod](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod_v1)
that conditionally enables secrets via Hashicorp Vault.

- If `var.enable_vault` is `true`
  - Enable Vault annotations
- If `var.enable_vault` is `false`
  - Disable Vault annotations
  - Mount a volume
  - Set an environment variable

```terraform
variable "enable_vault" {
  type = bool
}

resource "kubernetes_pod_v1" "test" {
  metadata {
    name = "terraform-example"
    annotations = var.enable_vault ? {
      "vault.hashicorp.com/client-timeout"            = "5m"
      "vault.hashicorp.com/agent-inject"              = "true"
      "vault.hashicorp.com/agent-inject-secret-azure" = "azure/creds/some-secret-name"
      "vault.hashicorp.com/role"                      = "some-role-name"
      "vault.hashicorp.com/agent-pre-populate-only"   = "true"
    } : {}
  }

  spec {
    container {
      image = "busybox"
      name  = "app"

      dynamic "env" {
        # A bogus map for a conditional block
        for_each = merge(var.enable_vault ? {} : { vault_disabled = true })

        content {
          name  = "DOCKER_CONFIG"
          value = "/dockerconfig"
        }
      }

      dynamic "volume_mount" {
        # A bogus map for a conditional block
        for_each = merge(var.enable_vault ? {} : { vault_disabled = true })

        content {
          name       = "dockerconfig"
          sub_path   = ".dockerconfigjson"
          mount_path = "/dockerconfig/config.json"
        }
      }
    }

    dynamic "volume" {
      # A bogus map for a conditional block
      for_each = merge(var.enable_vault ? {} : { vault_disabled = true })

      content {
        name = "dockerconfig"
        secret {
          default_mode = "0640"
          optional     = false
          secret_name  = "container-registry"
        }
      }
    }
  }
}
```

## References
- <https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/pod_v1>
- <https://www.terraform.io/language/expressions/dynamic-blocks>
- <https://www.terraform.io/language/functions/merge>
