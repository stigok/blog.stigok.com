---
layout: post
title: "dnsmasq dhcp not handing out addresses"
date: 2017-05-22 18:42:46 +0200
categories: dnsmasq dhcp
redirect_from:
  - /post/dnsmasq-dhcp-not-handing-out-addresses
---

Make sure the interface you are running on is configured with the same address as the `listen-address` in your `dnsmasq.conf`.

Configuring with `ip`

    sudo ip addr add dev enp3s0 192.168.70.1/24

And my `dnsmasq` config as follow

    keep-in-foreground
    
    # Disable DNS
    port=0
    
    # Setup DHCP
    log-dhcp
    interface=enp3s0
    listen-address=192.168.70.1
    dhcp-range=192.168.70.2,192.168.70.255,72h
    dhcp-option=option:router,192.168.70.1

To see what happens on the interface, `tcpdump` is always helpful

    # tcpdump --interface=enp3s0

## References
- `man dnsmasq`