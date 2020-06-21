---
layout: post
title:  "Setup a browser kiosk in NixOS with Xorg and Openbox"
date:   2020-06-20 01:53:35 +0200
categories: nixos xorg openbox kiosk
---

I wanted to create a kiosk-like setup for a box hooked up to a touch screen
TV. A problem I have yet to solve is to allow the users to configure wifi
through the web browser, but that is a problem for another post.

Note that this is not hardened for security, and it is possible to escape the
kiosk and get to the openbox desktop. However, this suited my needs, as I
want the users to be able to browse arbitrary websites and save files to the
disk.

## Prerequisites

- A machine with [NixOS][https://nixos.org/] installed
- **Basic** know-how of how to configure NixOS

## Setup

I have done my best to comment the file to make it easily understandable.
If you read the source, you should be able to understand what I am settings up
here.

Save this file to */etc/nixos/xorg.conf*.

```nix
{ pkgs, ... }:

let
  kioskUsername = "kiosk";
  browser = pkgs.firefox;
  autostart = ''
    #!${pkgs.bash}/bin/bash
    # End all lines with '&' to not halt startup script execution

    # https://developer.mozilla.org/en-US/docs/Mozilla/Command_Line_Options
    firefox --kiosk https://stigok.com/ &
  '';

  inherit (pkgs) writeScript;
in {
  # Set up kiosk user
  users.users = {
    "${kioskUsername}" = {
      group = kioskUsername;
      isNormalUser = true;
      packages = [ browser ];
    };
  };
  users.groups."${kioskUsername}" = {};

  # Configure X11
  services.xserver = {
    enable = true;
    layout = "us"; # keyboard layout
    libinput.enable = true;

    # Let lightdm handle autologin
    displayManager.lightdm = {
      enable = true;
      autoLogin = {
        enable = true;
        timeout = 0;
        user = kioskUsername;
      };
    };

    # Start openbox after autologin
    windowManager.openbox.enable = true;
    displayManager.defaultSession = "none+openbox";
  };

  # Overlay to set custom autostart script for openbox
  nixpkgs.overlays = with pkgs; [
    (self: super: {
      openbox = super.openbox.overrideAttrs (oldAttrs: rec {
        postFixup = ''
          ln -sf /etc/openbox/autostart $out/etc/xdg/openbox/autostart
        '';
      });
    })
  ];

  # By defining the script source outside of the overlay, we don't have to
  # rebuild the package every time we change the startup script.
  environment.etc."openbox/autostart".source = writeScript "autostart" autostart;
}
```

Import the xorg configuration in */etc/nixos/configuration.nix*

```nix

  imports =
    [ # Include the results of the hardware scan.
      ./hardware-configuration.nix
      ./xorg.nix
    ];

```

Now you can rebuild your configuration

```
# nixos-rebuild test
```

If you are making changes to the autostart script, you can do so without
having to wait for `openbox` to be rebuilt again, thanks to defining it
in a static file and linking it into the `openbox` output directory.

Note that the `display-manager.service` will not automatically restart
after rebuild. You will have to restart it manually.

```
# nixos-rebuild test
# systemctl restart display-manager.service
```

## References

https://github.com/NixOS/nixpkgs/pull/11273
http://openbox.org/wiki/Help:Autostart
https://wiki.archlinux.org/index.php/LightDM
https://wiki.archlinux.org/index.php/Openbox#Autostart

