---
layout: post
title: "Setting up wireless network on a Raspberry Pi 3 with Arch Linux Arm armv7h"
date: 2017-06-05 23:14:36 +0200
categories: systemd archarm linux networking wpa_supplicant
redirect_from:
  - /post/setting-up-wireless-network-on-a-raspberry-pi-3-with-arch-linux-arm-armv7h
---

There's a built-in wireless network card in the Raspberry Pi 3. I had some initial problems
getting the wireless card to work at first, but eventually it seemed to be mainly due to
invalid WPA key credentials or non-existing SSIDs.

The most helpful debugging tool here was `wpa_supplicant -c/etc/wpa_supplicant/wpa_supplicant-wlan0.conf -iwlan0 -d`
which put `wpa_supplicant` in foreground with debug mode, outputting to stdout.

I ended up managing the device through systemd, which was a first. Usually been doing it with ifupdown/interfaces or NetworkManager, but I liked this approach for a somewhat static configuration.

## Let systemd manage wlan card

Create a new file `/etc/systemd/network/wlan0.network`:

    [Match]
    Name=wlan0

    [Network]
    Description=On-board wireless NIC
    DHCP=yes

In the same directory, there probably already is a `eth0.network` file which you can look at
for quick reference, but all the information you need is in `man systemd.network`. To refresh
these settings, restart the service (probably is a signal for this, but this works too):

    # systemctl restart systemd-networkd

The device should now show up as *managed* and maybe even *configuring*, even though nothing
is configuring it at the moment.

    IDX LINK             TYPE               OPERATIONAL SETUP     
      1 lo               loopback           carrier     unmanaged 
      2 eth0             ether              routable    configured
      3 wlan0            wlan               no-carrier  configuring...

## Set up wpa_supplicant to handle wireless connections

Looking at the source of the pre-installed service file for interface specific configuration at
`/lib/systemd/system/wpa_supplicant@.service`, it expects wpa_supplicant configuration files to
reside at `/etc/wpa_supplicant/wpa_supplicant-%I.conf`, where `%I` is the interface name, i.e. `wlan0`.
This service file will be invoked when we later run it with `systemctl start wpa_supplicant@wlan0`.
If you try this now, it will probably fail with error stating that the *-wlan0.conf* file doesn't exist.

Using example configuration in `man wpa_supplicant.conf` and updating the settings to match my local
wireless access point (AP) settings. Now save this file where the wpa_supplicant service expects it to be;
`/etc/wpa_supplicant/wpa_supplicant-wlan0.conf`:

    ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=wheel
    network={
         ssid="sshowfosho"
         scan_ssid=1
         key_mgmt=WPA-PSK
         psk="yoforeal"
    }

Now start the service

    # systemctl start wpa_supplicant@wlan0.service

Check the journal for messages

    # journalctl -f

See if the device is getting an IP address

    # ip addr show wlan0

If this doesn't work, check the debug logs of wpa_supplicant as described at the top.
Drop comments if you are having troubles :)