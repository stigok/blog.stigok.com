---
layout: post
title:  "Use Raspberry Pi as WiFi AP and route traffic through Wireguard (port 53)"
date:   2019-03-26 19:48:54 +0100
categories: wireguard raspberrypi
---

## Introduction

I was at a place that was blocking traffic on all ports except 53 (DNS). That
had me thinking I could set up a WireGuard connection to tunnel traffic through
there. But I wanted more devices to be able to access it simultaneously, so I
set up a Raspberry Pi as a wireless access point and routed all the WiFi traffic
through the WireGuard tunnel.

- It is expected that you already have a [WireGuard server set up][1]
- [Configure a WireGuard server][1] interface to listen to port 53 (just set `ListenPort=53`)

This guide will result in the following network configuration:

#### Wireless interface:

- Static IP: `10.13.37.1/24`, `fd13:37::1/120`
- DHCP range: `10.13.37.0/24`, `fd13:37::/120`
- DNS server: `10.13.37.1`

#### Wireguard interface:

- Client: `10.8.0.101/24`, `fd00:8::101`
- Server: `10.8.0.1/24`, `fd00:8::1`
- DNS server: `91.239.100.100`
- Endpoint: `<wireguard server ip>:53`

## Setup steps
### Spoof ethernet MAC address

I don't want the other end of the line to know the original MAC of my network
device, so I changed it inside */boot/config.cmd* by appending
`smsc95xx.macaddr=XX:XX:XX:XX:XX:XX` to the end of the first line.

### Setup WireGuard client

Install prerequisites

```
$ sudo apt install raspberrypi-kernel-headers dirmngr
$ echo 'deb http://deb.debian.org/debian/ unstable main' | sudo tee -a /etc/apt/sources.list.d/unstable.list
$ sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 8B48AD6246925553
$ sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 7638D0442B90D010
$ sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 04EE7237B7D453EC
$ printf 'Package: *\nPin: release a=unstable\nPin-Priority: 150\n' | sudo tee --append /etc/apt/preferences.d/limit-unstable
$ sudo apt install wireguard
```

Create a new private key and output the public key by running the below command

```
$ sudo sh -c 'umask 077; wg genkey | tee /etc/wireguard/private.key | wg pubkey'
K1VaBTs6+09vmSjpXjDxDecMuTwDUZV6i3zf1u0kCXo=
```

[Configure your WireGuard server][1] to allow this peer by registering the
public key above. How to configure server side is not described in this guide.

Create a client configuration file *similiar* to the following, but with your
own specific modifications inside */etc/wireguard/wg0.conf*:

```
[Interface]
Address = 10.8.0.101/24, fd00:8::101/48
PrivateKey = <private key contents>
DNS = 91.239.100.100
PostUp   = netfilter-persistent start
PostDown = netfilter-persistent flush

[Peer]
PublicKey = <public key of server>
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = <server ip>:53
Endpoint = <server ipv6>:53
PersistentKeepalive = 15
```

Auto-start the interface on boot

```
$ sudo systemctl enable wg-quick@wg0
```

Reboot (note that the Pi will now probably have a different IP due to the change of the MAC address)

```
$ sudo reboot
```

Verify that your WireGuard interface connects successfully before continuing.
It will probably make troubleshooting easier in case of problems.

### Setup as WiFi access point (AP)

Install required programs, and stop them for now. They will start automatically
on boot anyway.

```
$ sudo apt install dnsmasq hostapd
$ sudo systemctl stop dnsmasq hostapd
```

Set a static IP address for the wireless interface.
Append the following to */etc/dhcpcd.conf*:

```
interface wlan0
    static ip_address=10.13.37.1/24
    static ip6_address=fd13:37::1/120
    nohook wpa_supplicant
```

Configure the DHCP server by replacing **all the contents** of */etc/dnsmasq.conf*
with the following:

```
interface=wlan0
except-interface=eth0
except-interface=wg0
listen-address=10.13.37.1
dhcp-range=10.13.37.2,10.13.37.254,255.255.255.0,24h
dhcp-option=option:dns-server,10.13.37.1
listen-address=fd13:37::01
dhcp-range=fd13:37::02,fd13:37::ff,24h
dhcp-authoritative
enable-ra
bogus-priv
domain-needed
```

Configure the WiFi AP server deamon by pasting the following config in
*/etc/hostapd/hostapd.conf*. Take note of the `ssid` and `wpa_passphrase`
values and change them to your desire.

```
interface=wlan0
driver=nl80211
ssid=NETWORK_NAME
hw_mode=g
channel=7
wmm_enabled=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=WIFI_PASSPHRASE
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP
```

Specify the location of the previous configuration file inside
*/etc/default/hostapd* with the following value:

```
DAEMON_CONF="/etc/hostapd/hostapd.conf"
```

Enable IPv4 forwarding

```
$ echo 'net.ipv4.ip_forward=1' | sudo tee /etc/sysctl.d/97-wifi-ap.conf
```

Masquerade outbound traffic on wg0

```
$ sudo iptables -A FORWARD -i wlan0 -o wg0 -j ACCEPT
$ sudo iptables -A FORWARD -i wg0 -o wlan0 -m state --state RELATED,ESTABLISHED -j ACCEPT
$ sudo iptables -t nat -A  POSTROUTING -o wg0 -j MASQUERADE
```

Persist iptables rules by installing `iptables-persistent` and answer yes
in the prompt. (If the prompt fails or you for other reasons want to do
this again later, use `dpkg-reconfigure iptables-persistent`, or manually
run `netfilter-persistent`).

```
$ sudo apt install iptables-persistent
```

Unmask and enable services on boot

```
$ sudo systemctl unmask hostapd
$ sudo systemctl enable hostapd dnsmasq
```

Reboot once again, and see if it works as expected. You should be able to connect
to the Pi through the WiFi AP. Once connected, verify that you have an IP address
different to that of which you started with.

```
$ curl ip.stigok.com
```

## References
- https://pimylifeup.com/raspberry-pi-mac-address-spoofing/
- https://www.raspberrypi.org/documentation/configuration/wireless/access-point.md
- https://github.com/adrianmihalko/raspberrypiwireguard
- https://www.ckn.io/blog/2017/12/28/wireguard-vpn-portable-raspberry-pi-setup/

[1]: https://blog.stigok.com/2018/10/08/wireguard-vpn-server-on-centos-7.html
