---
layout: post
title:  "Prepare a headless Raspberry Pi 3 with Arch Linux ARM to connect to WiFi"
date:   2019-09-03 00:58:32 +0200
categories: raspberrypi arch
---

## Prepare SD-card

Prepare the SD-card to use for the Raspberry Pi 3 by following the
[official documentation for Arch Linux ARM](https://archlinuxarm.org/platforms/armv8/broadcom/raspberry-pi-3#installation).

## Prepare headless operation

It is expected that the root file system on the SD-card is mounted on your
system at */mnt/root*.

**Beware** that both *root* and *alarm* user has default passwords, hence the
Pi should never be put on an untrusted network before those are changed
(there's a reason the username is *alarm*).

Create a symlink to the `wpa_supplicant` systemd instance service file to make
it start on boot.

```
# cd /mnt/root/etc/systemd/system/multi-user.target.wants
# ln -s /usr/lib/systemd/system/wpa_supplicant@.service wpa_supplicant@wlan0.service
```

`dhcpcd` also has to start at boot in order to handle DHCP for the wireless
interface when it connects

```
# ln -s /usr/lib/systemd/system/dhcpcd.service
```

Create a `wpa_supplicant` configuration file with your wifi connection details

```
# wpa_passphrase "My SSID" "My passphrase" > /mnt/root/etc/wpa_supplicant/wpa_supplicant-wlan0.conf
```

Place your public key inside the SSH user configuration directory to allow for
root login

```
# mkdir /mnt/root/root/.ssh
# chmod 0700 /mnt/root/root/.ssh
# cat ~/.ssh/id_rsa.pub > /mnt/root/root/.ssh/authorized_keys
# chmod 0600 /mnt/root/root/.ssh/authorized_keys
```

Boot up the Pi and it should connect to the wireless network you specified
automatically on boot.

When I'm unsure what IP a device on my network has got, I use `arp-scan` to
probe.

```
# arp-scan -l
Interface: wlp3s0, datalink type: EN10MB (Ethernet)
Starting arp-scan 1.9.5 with 256 hosts (https://github.com/royhills/arp-scan)
192.168.1.1	xx:xx:xx:xx:xx:a6	xxx
192.168.1.22	xx:xx:xx:xx:xx:24	xxx
192.168.1.26	xx:xx:xx:xx:xx:a4	Raspberry Pi Foundation

3 packets received by filter, 0 packets dropped by kernel
Ending arp-scan 1.9.5: 256 hosts scanned in 2.013 seconds (127.17 hosts/sec). 3 responded
```

Then attempt SSH

```
$ ssh root@192.168.1.26 -i ~/.ssh/id_rsa
Last login: Mon Jul 22 14:30:50 2019 from 192.168.1.19
[root@alarm ~]#
```

Apparently, I need to get `ntpd` or `systemd-timesyncd` running.
