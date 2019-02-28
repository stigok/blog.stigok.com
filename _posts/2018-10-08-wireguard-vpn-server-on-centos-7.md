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

Configure the wireguard network interface

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

Take note of the public key which all of the clients will need in order to
establish a wireguard connection to this server

```terminal
# wg
interface: wg0-server
  public key: rWIm1TmlDowftuzcp0uipzJQ5FCzTc2brMLzYmirZXc=
  private key: (hidden)
  listening port: 34777
```

## Client

**This host is running Arch Linux**

Install wireguard packages

```terminal
# pacman -S linux-headers wireguard-dkms wireguard-utils
```

Create a folder only accessible by root and generate a private key

```terminal
# mkdir /etc/wireguard/
# chmod 0600 /etc/wireguard
# wg genkey > /etc/wireguard/private.key
```

Create a `wg-quick` configuration file which makes it easier to bring up and
down one or more WireGuard interfaces

```terminal
# tee /etc/wireguard/wg0.conf
[Interface]
Address = 10.7.0.2/24, fd00:7::2/48
PrivateKey = <replace with the contents of your private key>
DNS = 10.7.0.1

[Peer]
PublicKey = <replace with server public key string>
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = <replace with server IPv6 address/hostname>:34777
Endpoint = <replace with server IPv4 address/hostname>:34777
PersistentKeepalive = 15
```

- Remove IPv6 configuration if you're not using it
- Replace the `Address` with the IP address intended for this client
- Replace the value of `PrivateKey` with the contents of the *private.key* file
- Replace the value of `PublicKey` with the public key of the server that was
  determined in the first section
- Replace the `Endpoint` with the IP or hostname of the server

## **Server** firewall configuration

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
# firewall-cmd --add-service=wireguard zone=public --permanent
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

Enable IPv4 forwarding and make it persistent across reboots

```terminal
# sysctl net.ipv4.ip_forward=1
# tee -a /etc/sysctl.d/99-sysctl.conf
sysctl net.ipv4.ip_forward=1
```

Start server client

```terminal
# systemctl start wg-quick@wg0-server.service
```

## Client

Attempt to connect to the server

```terminal
# wg-quick up wg0
```

See if you are able to send and receive traffic, and at the same time check
your IP address

```terminal
$ curl -i -4 ip.stigok.com
```

## References
- https://www.wireguard.com/install/
- https://fedoraproject.org/wiki/EPEL/FAQ
- https://git.zx2c4.com/WireGuard/about/src/tools/man/wg-quick.8
- https://serverfault.com/questions/664576/should-iptables-be-this-long-many-chains
- https://www.digitalocean.com/community/tutorials/how-to-set-up-a-firewall-using-firewalld-on-centos-7

[epel-faq]: https://fedoraproject.org/wiki/EPEL/FAQ#What_is_EPEL.3F
