---
layout: post
title:  "How to make wpa_cli work in NixOS"
date:   2021-05-04 22:31:28 +0200
categories: nixos
excerpt: wpa_cli will not work out of the box in NixOS with wpa_supplicant. Here's the missing config line.
#proccessors: pymd
---

## Preface

More than once I've found myself stuck at a place without a network connection
in need of scanning for nearby access points. However, `wpa_cli` does not work
out of the box when you've only set up your known SSID's and pre-shared keys using

```nix
networking.wireless = {
  enabled = true;
  networks."my-ssid".psk = "helloworld42";
};
```

## Enabling the control interface

The following line tends to be put at the top of the `wpa_supplicant` configuration
file:

```
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=wheel
```

This makes the control interface accessible by the `wheel` group.
The manpage for wpa_supplicant.conf(5) states

> [ctrl_interface]
> allow frontend (e.g., wpa_cli) to be used by all users in 'wheel' group

I've always added this line when I configured
my `/etc/wpa_supplicant.conf` by hand in other operating systems, but I guess I forgot why I did it.

In NixOS this line should be added using `extraConfig`:

```nix
networking.wireless.extraConfig = ''
  ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=wheel
'';
```

The resulting wpa_supplicant.conf will then look something (if not exactly) like this:

```
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=wheel

network={
  ssid="my-ssid"


  psk=fec51ab9d2363e43a2f5e454aa3eab77da1aa3ae21ba71ee806e1e1f5d3cf7bd

}
```


## References
- [man 5 wpa_supplicant.conf](https://linux.die.net/man/5/wpa_supplicant.conf)
- [clever at #nixos @ freenode.net](https://logs.nix.samueldr.com/nixos/2020-04-30#3391000)
