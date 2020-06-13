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

- <https://nixos.org/nix/manual/>
  - [Language constructs](https://nixos.org/nix/manual/#sec-constructs)
- <https://nixos.org/nixos/manual/>
- <https://nixos.org/nixpkgs/manual/>
  - [Library functions](https://nixos.org/nixpkgs/manual/#sec-functions-library)

## Package management

### Search for package

```
$ nix search packagename
```

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

### Fetch git repository from private server over SSH

`pkgs.fetchgit` does not use the local *ssh_config*. Use `builtins.fetchGit`
instead.

```
src = builtins.fetchGit {
  url = "ssh://git@git.stigok.com/stigok/utils.git";
  ref = "master";
  rev = "f8bdc053406ad28ef4b6cbb29e418ce69f31f05f";
};
```

Reference: https://nixos.org/nix/manual/#ssec-builtins

### Change package source to subdirectory after fetchGit

I have a monorepo with multiple utility programs inside. I have to change
root to a subdirectory, using `sourceRoot`, after I've downloaded its sources.

```
{ pkgs ? import <nixpkgs> {} }:
  let
    package = (import ./default.nix { inherit pkgs system; }).package;
    newSrc  = builtins.fetchGit {
      url = "ssh://git@git.stigok.com/stigok/utils.git";
      ref = "master";
      rev = "f8bdc053406ad28ef4b6cbb29e418ce69f31f05f";
    };
  in
    package.overrideAttrs (old: {
      src = newSrc;
      sourceRoot = "${newSrc.outPath}/{package.packageName}";
    })
```

The original package derivation resides in
*./default.nix* and the package tree that `fetchGit` downloads looks like this:

```
README.md
my-package-src-in-a-subdir/
my-package-src-in-a-subdir/package.json
my-package-src-in-a-subdir/package-lock.json
my-package-src-in-a-subdir/index.js
```

I can now build the overridden package using `nix-build override.nix`.

References:
- https://github.com/tfc/node2nix_bootstrap/blob/master/override.nix
- infinisil #nixos @ Freenode

### Create a tar archive of a fetchGit

The title could have been "create a tar archive", but using it straight after
`builtins.fetchGit` is what I needed to know.

```
gitSrc = builtins.fetchGit {
  url = "ssh://git@git.stigok.com/stigok/utils.git";
  ref = "master";
  rev = "c3c4be782efd79518ddbac006a7fcb3532c35f7e";
};
tarSrc = runCommand "tar-src" {} "${gnutar}/bin/tar -cf $out ${gitSrc}";
```

`tarSrc` now contains a path to the tar file.

References:
- euank #nixos @ Freenode
- https://nixos.org/nixpkgs/manual/#chap-trivial-builders

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
