---
layout: post
title:  "My (expanding) list of usage notes for Nix and NixOS"
date:   2020-03-07 19:18:27 +0100
categories: nixos draft
---

### String interpolation
### Use variable as map/object key

```nix
let
  customUser = "stigok";
in {
  users.extraUsers = {
    ${customUser} = {
      systemUser = true;
      description = "This is a user called ${customUser}";
    };
  };
}
```

### Check a file for syntax errors

```shell
$ nix-instantiate --parse <file>
$ nix-instantiate --parse myfile.nix
```

### Pass function arguments with nix-instantiate

```shell
$ nix-instantiate --arg config <val1> --arg nixpkgs <val2> myfile.nix
$ nix-instantiate --arg config '{}' --arg pkgs '<nixpkgs>' --eval postgresql.nix
```

Reference: slack1256 #nixos @ Freenode

## References

- <https://nixery.dev/nix-1p.html>
