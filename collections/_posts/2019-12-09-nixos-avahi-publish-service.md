---
layout: post
title:  "Publish a service with Avahi on NixOS"
date:   2019-12-09 19:01:12 +0100
categories: nixos dns
---

Usually I would create a new *.service* file and fill in the service details there.
But for now I was just able to figure out how to broadcast the hostname and IP
address of the server the NixOS way.

I should later find out how to publish a HTTP web service listening on a specific port.
This post will in that case be updated. Until then -- this will suffice!

## Configuration

Open up the main NixOS configuration file at */etc/nixos/configuration.nix* and
enable the avahi service:

```nix
# Publish this server and its address on the network
services.avahi = {
  enable = true;
  publish = {
    enable = true;
    addresses = true;
    workstation = true;
  };
};
```

Rebuild NixOS

```
# nixos-rebuild switch --fast
```

Now, on my other machine, I can see the server broadcasting itself

```
# systemctl start avahi-daemon.service
$ avahi-browse -park | grep nix
+;wlp3s0;IPv4;nix\032\09100\05830\05818\058a8\058e9\05850\093;_workstation._tcp;local
=;wlp3s0;IPv4;nix\032\09100\05830\05818\058a8\058e9\05850\093;_workstation._tcp;local;nix.local;192.168.0.2;9;
```


## References
- <https://jarmac.org/category/nixos.html>
- <https://github.com/NixOS/nixpkgs/blob/master/nixos/modules/services/networking/avahi-daemon.nix>
- <https://wiki.archlinux.org/index.php/Avahi#Adding_services>

[avahi-daemon.nix]: https://github.com/NixOS/nixpkgs/blob/master/nixos/modules/services/networking/avahi-daemon.nix
