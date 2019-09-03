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

Let's let `dhcpcd` take care of wireless networking (it's not handling wired
connections with this setup):

```
# ln -s /usr/lib/systemd/system/dhcpcd@.service /mnt/root/etc/systemd/system/multi-user.target.wants/dhcpcd@wlan0.service
```

Add a hook so that `dhcpcd` can take care of starting `wpa_supplicant`:

```
# ln -s /usr/share/dhcpcd/hooks/10-wpa_supplicant /mnt/root/usr/lib/dhcpcd/dhcpcd-hooks/
```

Create a `wpa_supplicant` configuration file with your wifi connection details

```
# echo "ctrl_interface=DIR=/var/run/wpa_supplicant" > /mnt/root/etc/wpa_supplicant/wpa_supplicant-wlan0.conf
# wpa_passphrase "My SSID" "My passphrase" >> /mnt/root/etc/wpa_supplicant/wpa_supplicant-wlan0.conf
```

Place your public key inside root user's SSH configuration directory to allow for
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

## Troubleshooting
### Network time synchronisation

I had troubles using `systemd-timesyncd`, as it was failing all the time.
And, initially, I had problems with the default config of `openntpd`, but got
it working by avoiding hostnames for the `server` directives, and instead
using plain IPv4 addresses. I got the IPs of the DNS records of `pool.ntp.org`
(`dig +short pool.ntp.org`).

`/etc/ntpd.conf` (on the Pi):

```
server 194.192.112.20
server 92.246.24.228
server 193.162.159.97
server 5.103.128.88
```

Add a hook at `/etc/dhcpcd.exit-hook` (on the Pi):

```
#!/bin/bash

if $if_up; then
	# Don't wait for this to finish
	systemctl start openntpd.service &
elif $if_down; then
	systemctl stop openntpd.service
fi
```
