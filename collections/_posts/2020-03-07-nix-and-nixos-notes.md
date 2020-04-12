---
layout: post
title:  "My (expanding) list of usage notes for Nix, NixOS and NixOps"
date:   2020-03-07 19:18:27 +0100
categories: nixos draft
---

## Manual page sources

When it comes to Nix, there are some different manuals which overlaps in terms
of , however, they don't contain the same amount of information. In order to
find all information available for a certain topic you might look for, prepare
to look through all of the following sources:

- https://nixos.org/nix/manual/
- https://nixos.org/nixos/manual/
- https://nixos.org/nixpkgs/manual/
  - [Library functions](https://nixos.org/nixpkgs/manual/#sec-functions-library)

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

### Inject `config` and `services` variables to the nix repl

This has to be run on a machine running NixOS itself

```shell
$ nix repl '<nixpkgs/nixos>'
> config
> config.services
```

Reference: energizer bqv[m] clever #nixos @ Freenode

### Install a package for current user

```
$ nix-env -iA nixos.thunderbird
```

This will make the package available for the current logged in user only, *except*
when logged in as root, which will make it available for everyone.

Reference: https://nixos.org/nixos/manual/index.html#sec-ad-hoc-packages

## NixOps

### switch-to-configuration throws error deployment fails

```
[...]
trivial> /nix/var/nix/profiles/system/bin/switch-to-configuration: line 3: use: command not found
[...]
```

I was getting erros while attempting to deploy to a specific machine in my
NixOps deployment. It's a 32-bit Jetway box, and so it needs 32-bit packages.
When deploying from a different architecture, like in my case, a 64-bit intel,
nixops needs to know what system it is targetting explicitly.

```
{
  network.enableRollback = true;
  network.description = "private infra";
  nix =
    { resources, ... }:
    {
      imports = [
        ./servers/nix/configuration.nix
      ];
      deployment.targetHost = "192.168.0.2";
      nixpkgs.system = "i686-linux";
    };
}
```

Reference: <https://github.com/NixOS/nixops/issues/864>

### Determine NixOS machine architecture

To figure out what architecture your box is running, you can run `nix-info` or
for example `nix-eval`:

```nix
$ nix-info
system: "x86_64-linux", multi-user?: yes, version: nix-env (Nix) 2.3.3, channels(username): "", channels(root): "nixos-19.09.2213.71c6a1c4a83", nixpkgs: /nix/var/nix/profiles/per-user/root/channels/nixos

$ nix eval nixpkgs.system
"x86_64-linux"
```

## References

- <https://nixery.dev/nix-1p.html>
