---
layout: post
title:  "Setting up a WireGuard VPN server on CentOS"
date:   2018-10-08 11:01:47 +0200
categories: wireguard centos archlinux
---

## Introduction

I'm tired of OpenVPN quirks and configuration issues across my devices.
Additionally, I've been planning to try out WireGuard for some time now, after
hearing praises from some different people along my way.

I set up a new CentOS box to act as the VPN server, and the client in my
guide is, as usual, running Arch Linux.

The section headers below tells you whether the work is on the CentOS server
(server) or the Arch Linux client (client).

## Server

**This host is running CentOS**

> The installation steps are based on <https://www.wireguard.com/install/> and
> may have been updated since time of writing this post.

Add the [EPEL][epel-faq] (Extra Packages for Enterprise Linux) RPM repo and
install WireGuard and utilities

```terminal
# curl -Lo /etc/yum.repos.d/wireguard.repo https://copr.fedorainfracloud.org/coprs/jdoss/wireguard/repo/epel-7/jdoss-wireguard-epel-7.repo
# yum install epel-release
# yum install wireguard-dkms wireguard-tools
```

Create an empty server config file with proper permissions

```terminal
# mkdir /etc/wireguard && cd /etc/wireguard
# bash -c 'umask 077; touch wg0-server.conf'
```

Configure the wireguard network interface. Here we are using the output of `wg genkey` directly. The `PrivateKey` option in the `wg-quick` configuration file also accepts a file path to a file containing the private key, if that should be more desirable.

```terminal
# ip link add dev wg0-server type wireguard
# ip addr add dev wg0-server 10.7.0.1/24
# ip addr add dev wg0-server fd00:7::1/48
# wg set wg0-server listen-port 34777 private-key <(wg genkey)
```

Save configuration to a file

```terminal
# wg-quick save wg0-server
```

Take note of the public key of the server. All of the clients will need it in order to
establish a wireguard connection to this server.

```terminal
# wg
interface: wg0-server
  public key: a0ap6Ze3Ug9OXNhd+w6xAj4gawL2b//uZsVab1ToJAg=
  private key: (hidden)
  listening port: 34777
```

## Client

This example client host is running Arch Linux. If you are running CentOS on your client too, repeat the installation steps as described in the previous step instead.

Install wireguard packages.

```terminal
# pacman -S linux-headers wireguard-dkms wireguard-utils
```

Create a folder only accessible by root and generate a private key

```terminal
# mkdir /etc/wireguard/
# chmod 0600 /etc/wireguard
# wg genkey > /etc/wireguard/private.key
```

The contents of `private.key` will now look something like this:

```terminal
# cat /etc/wireguard/private.key
oMGknnGRL0MO9gmLRNNG8H+2yGHSroo2r6U95WfchHQ=
```

This is a **private key**, hence, a secret that should not be shared with anyone.
In contrast to the **public key** which is not considered secret and can even
be sent over an unencrypted channel.

The public key has to be registered on the server in a later step. Take note of how
to extract it:

```terminal
# wg pubkey < /etc/wireguard/private.key
BwVtKNSF50L973aQ/YT+s/3lmlLjcbhkwp4uELqwEVU=
```

Create a `wg-quick` configuration file which makes it easier to bring up and
down one or more WireGuard interfaces

```
# tee /etc/wireguard/wg0.conf
[Interface]
Address = 10.7.0.2/24, fd00:7::2/48
PrivateKey = oMGknnGRL0MO9gmLRNNG8H+2yGHSroo2r6U95WfchHQ=

[Peer]
PublicKey = a0ap6Ze3Ug9OXNhd+w6xAj4gawL2b//uZsVab1ToJAg=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 2001:db8:600d:deed::7:34777
Endpoint = 203.0.113.77:34777
PersistentKeepalive = 15
```

- Remove IPv6 configuration if you're not using it.
- Replace the `Address` with the IP address intended for this client, and the
  prefix for the subnet it belongs to.
- Replace the value of `PrivateKey` with the contents of the client *private.key* file.
- Replace the value of `PublicKey` with the public key of the server that was
  determined in a previous section.
- Replace the `Endpoint` with the public IP or hostname of the server.

## **Server** firewall configuration

> If you're not running `firewalld`, this step may be skipped.

Create a new firewalld service definition for WireGuard

```terminal
# tee /etc/firewalld/services/wireguard.xml
<?xml version="1.0" encoding="utf-8"?>
<service>
  <short>wireguard</short>
  <description>WireGuard (wg) custom installation</description>
  <port protocol="udp" port="34777"/>
</service>
```

Enable the custom WireGuard service in firewalld

```terminal
# firewall-cmd --add-service=wireguard --zone=public --permanent
```

Enable masquerading

```terminal
# firewall-cmd --zone=public --add-masquerade --permanent
```

Reload firewalld and take a look at the zone configuration

```terminal
# firewall-cmd --reload
# firewall-cmd --list-all
public (active)
  target: default
  icmp-block-inversion: no
  interfaces: eno1
  sources:
  services: ssh dhcpv6-client wireguard
  ports:
  protocols:
  masquerade: yes
  forward-ports:
  source-ports:
  icmp-blocks:
  rich rules:
```

Enable IPv4 forwarding, and, if applicable, IPv6 forwarding as well.

```terminal
# sysctl -w net.ipv4.ip_forward=1
# sysctl -w net.ipv6.conf.all.forwarding=1
```

Make sysctl settings persistent across reboots

```terminal
# tee -a /etc/sysctl.d/99-sysctl.conf
sysctl -w net.ipv4.ip_forward=1
sysctl -w net.ipv6.conf.all.forwarding=1
```

Start the wireguard client process on the server

```terminal
# systemctl start wg-quick@wg0-server.service
```

### Letting VPN clients connect to each other (optional)

Add wireguard interface to the `internal` firewalld zone

```terminal
# firewall-cmd --add-interface=wg0-server --zone=internal --permanent
```

Enable masquerading

```terminal
# firewall-cmd --zone=internal --add-masquerade --permanent
```

Enable the services you'd like to be available on this network using `firewall-cmd --zone=internal --add-service [name]`. Remember to add the flag `--permanent` when it works. So easy to forget.

## Server - Allow client to connect

Now, back in the *server* configuration file, add the client public key and the IP's it should be allowed to register with on the server.
Copy the public key from the client and paste it into the server configuration, like below.

```
[Peer]
PublicKey = BwVtKNSF50L973aQ/YT+s/3lmlLjcbhkwp4uELqwEVU=
AllowedIPs = 10.7.0.2/32
AllowedIPs = fd00:7::2/128
```

Restart the server. No errors should be thrown.

```terminal
# systemctl restart wireguard@wg0-server
```

## Client - Attempt to connect to server

Attempt to connect to the server

```terminal
# wg-quick up wg0
```

See if you are able to send and receive traffic, and at the same time check
your IP address

```terminal
$ curl -i -4 ip.stigok.com
```

## Complete configuration

### Server

```terminal
# cat /etc/wireguard/wg0-server.conf
[Interface]
Address = 10.7.0.1/24, fd00:7::1/48
ListenPort = 34777
PrivateKey = 0GKvQPZ4oT236J+PrWo5/OO67nwSxJ+/p7N+hBw3WHU=

[Peer]
PublicKey = BwVtKNSF50L973aQ/YT+s/3lmlLjcbhkwp4uELqwEVU=
AllowedIPs = 10.7.0.2/32
AllowedIPs = fd00:7::2/128
```

### Client

```terminal
# cat /etc/wireguard/private.key
oMGknnGRL0MO9gmLRNNG8H+2yGHSroo2r6U95WfchHQ=
```

```terminal
# cat /etc/wireguard/wg0.conf
[Interface]
Address = 10.7.0.2/24, fd00:7::2/48
PrivateKey = oMGknnGRL0MO9gmLRNNG8H+2yGHSroo2r6U95WfchHQ=
# Ability to specify DNS servers to be picked up by resolvconf
#DNS = 10.7.0.1, fd00:7::1

[Peer]
# This is the public key of the server
PublicKey = a0ap6Ze3Ug9OXNhd+w6xAj4gawL2b//uZsVab1ToJAg=
AllowedIPs = 0.0.0.0/0, ::/0
# These are example public addresses to reach the server
Endpoint = 2001:db8:600d:deed::7:34777
Endpoint = 203.0.113.77:34777
PersistentKeepalive = 15
```

## References
- https://www.wireguard.com/install/
- https://fedoraproject.org/wiki/EPEL/FAQ
- https://git.zx2c4.com/WireGuard/about/src/tools/man/wg-quick.8
- https://serverfault.com/questions/664576/should-iptables-be-this-long-many-chains
- https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7
- https://www.ripe.net/participate/member-support/lir-basics/ipv6_reference_card.pdf

[epel-faq]: https://fedoraproject.org/wiki/EPEL/FAQ#What_is_EPEL.3F
