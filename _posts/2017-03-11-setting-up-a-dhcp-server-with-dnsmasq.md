---
layout: post
title: "Setting up a dhcp server with dnsmasq"
date: 2017-03-11 23:09:03 +0100
categories: dnsmasq dhcp
redirect_from:
  - /post/setting-up-a-dhcp-server-with-dnsmasq
---

Using a minimal _dnsmasq-dhcp-only.conf_

    # Disable DNS
    port=0
    log-dhcp
    interface=enp2s1
    listen-address=192.168.70.1
    dhcp-range=192.168.70.2,192.168.70.255,1h
    dhcp-option=option:router,192.168.70.1

Make sure the interface you are using is on the same subnet as specified in the `dhcp-range`, otherwise you will get errors like below

    # dnsmasq -C ./dnsmasq-dhcp-only.conf -d
    dnsmasq: started, version 2.75 DNS disabled
    dnsmasq: compile time options: IPv6 GNU-getopt DBus i18n IDN DHCP DHCPv6 no-Lua TFTP conntrack ipset auth DNSSEC loop-detect inotify
    dnsmasq-dhcp: DHCP, IP range 192.168.70.2 -- 192.168.70.255, lease time 3d
    dnsmasq-dhcp: no address range available for DHCP request via enp2s1

If this is the case, configure your interface

    # ifconfig enp2s1 192.168.70.1 netmask 255.255.255.0

And retry

    # dnsmasq -C /etc/dnsmasq-dhcp-enp2s1.conf -d           
    dnsmasq: started, version 2.75 DNS disabled
    dnsmasq: compile time options: IPv6 GNU-getopt DBus i18n IDN DHCP DHCPv6 no-Lua TFTP conntrack ipset auth DNSSEC loop-detect inotify
    dnsmasq-dhcp: DHCP, IP range 192.168.70.2 -- 192.168.70.255, lease time 3d
    dnsmasq-dhcp: 2063584936 available DHCP range: 192.168.70.2 -- 192.168.70.255
    dnsmasq-dhcp: 2063584936 client provides name: p-bc42
    dnsmasq-dhcp: 2063584936 DHCPDISCOVER(enp2s1) 00:42:ye:aa:bb:cc 
    dnsmasq-dhcp: 2063584936 tags: enp2s1
    dnsmasq-dhcp: 2063584936 DHCPOFFER(enp2s1) 192.168.70.192 00:42:ye:aa:bb:cc 
    dnsmasq-dhcp: 2063584936 requested options: 1:netmask, 28:broadcast, 3:router, 6:dns-server, 
    dnsmasq-dhcp: 2063584936 requested options: 15:domain-name, 12:hostname
    dnsmasq-dhcp: 2063584936 next server: 192.168.70.1
    dnsmasq-dhcp: 2063584936 sent size:  1 option: 53 message-type  2
    dnsmasq-dhcp: 2063584936 sent size:  4 option: 54 server-identifier  192.168.70.1
    dnsmasq-dhcp: 2063584936 sent size:  4 option: 51 lease-time  3d
    dnsmasq-dhcp: 2063584936 sent size:  4 option: 58 T1  1d12h
    dnsmasq-dhcp: 2063584936 sent size:  4 option: 59 T2  2d15h
    dnsmasq-dhcp: 2063584936 sent size:  4 option:  1 netmask  255.255.255.0
    dnsmasq-dhcp: 2063584936 sent size:  4 option: 28 broadcast  192.168.70.255
    dnsmasq-dhcp: 2063584936 sent size:  4 option:  3 router  192.168.70.1