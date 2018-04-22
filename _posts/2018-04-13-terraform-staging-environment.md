---
layout: post
title: "Terraform staging environment"
date: 2018-04-13 18:49:20 +0200
categories: terraform
redirect_from:
  - /post/terraform-staging-environment
---

I was looking for a fairly simple way of handling different deployment environments with Terraform. Here are some of my notes on variables, workspaces, environments and variable overrides.

- Use workspaces to separate state files. This is important when using variable overrides to e.g. go from `staging` to `testing` by changing namespaces. Without a change of workspace, Terraform will think you want to destroy `staging`, then set up `testing, instead of having them side by side.
- All files that ends with `.tf` will be loaded and appended onto each other.
- Resources with the same name will raise validation errors.
- `terraform.tfvars` and any other files that ends with `.auto.tfvars` will be automatically loaded, overriding default variables.
- Other `.tfvars` files may be references when running commands, e.g. `terraform plan -var-file=staging.tfvars`

Variable files containing overrides (`.tfvars`) have a simple syntax:

    gce_region = "europe-west1-c"
    namespace = "myapp-staging"

## References
- https://groups.google.com/forum/#!topic/terraform-tool/l6FGol0iXww
- https://www.terraform.io/docs/configuration/variables.html