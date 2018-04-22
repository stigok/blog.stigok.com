---
layout: post
title: "Watch dump dhcp requests on local network"
date: 2017-11-30 23:22:42 +0100
categories: tcpdump network
redirect_from:
  - /post/watch-dump-dhcp-requests-on-local-network
---

I am already connected to the local network, and I'm about to plug in a new, headless, device on the network. Since there are several other devices on the network, it's not trivial to find the IP address of it using `arp-scan` alone.

Start dumping all DHCP traffic on the local network. Here, my network interface is `enp3s0`, and using `-n` to avoid hostname resolution and `-e` to also dump the link-layer address (MAC address) of the hosts.

First I'm starting `tcpdump`, then I'm plugging in the device

    # tcpdump -nei enp3s0 port 67 or port 68
    23:00:12.566761 9a:a6:02:00:00:00 > ff:ff:ff:ff:ff:ff, ethertype IPv4 (0x0800), length 381: 0.0.0.0.68 > 255.255.255.255.67: BOOTP/DHCP, Request from 9a:a6:02:00:00:00, length 339

> DHCP is using TCP ports 67 and 68

Right after plugging it in, traffic appears. Now, I can run `arp-scan` on the same interface to map it to an IP address. Filtering through `grep` to reduce amount of output.

    # arp-scan -lI enp3s0 | grep 9a:a6:02:00:00:00
    10.10.3.146	9a:a6:02:00:00:00	(Unknown)

## References
- `man tcpdump`
- `man arp-scan`
- http://ask.xmodulo.com/monitor-dhcp-traffic-command-line-linux.html