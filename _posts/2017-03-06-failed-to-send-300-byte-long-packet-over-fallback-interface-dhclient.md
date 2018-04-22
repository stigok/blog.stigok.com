---
layout: post
title: "Failed to send 300 byte long packet over fallback interface (dhclient)"
date: 2017-03-06 17:34:45 +0100
categories: dhcp linux
redirect_from:
  - /post/failed-to-send-300-byte-long-packet-over-fallback-interface-dhclient
---

dhclient started complaining a lot recently:

    dhclient[17182]: DHCPREQUEST on wlp2s0 to 10.10.3.1 port 67
    kernel: [UFW AUDIT] IN= OUT=wlp2s0 SRC=10.10.3.227 DST=10.10.3.1 LEN=328 TOS=0x00 PREC=0x00 TTL=64 ID=56217 DF PROTO=UDP SPT=68 DPT=67 LEN=308
    dhclient[17182]: send_packet: Operation not permitted
    dhclient[17182]: dhclient.c:2593: Failed to send 300 byte long packet over fallback interface.

I have seen the `send_packet: Operation not permitted` messages before, and they have usually showed themselves when I'm blocking something I shouldn't be blocking. Might have a little _too_ strict iptables rules. But blocking the DHCP port 67 is usually not something I'd want to do.

I am using UFW for handling my iptables rules. As one of my first rules are
`ufw deny out on wlp2s0`, I need to place this one before it in the chain.

    # ufw insert 1 allow out on wlp2s0 to any port 67 proto udp comment DHCP

Now my dhclient seems to be happy again

    dhclient[17182]: DHCPREQUEST on wlp2s0 to 10.10.3.1 port 67
    kernel: [UFW AUDIT] IN= OUT=wlp2s0 SRC=10.10.3.227 DST=10.10.3.1 LEN=328 TOS=0x00 PREC=0x00 TTL=64 ID=56842 DF PROTO=UDP SPT=68 DPT=67 LEN=308
    dhclient[17182]: DHCPACK from 10.10.3.1
    dhclient[17182]: bound to 10.10.3.227 -- renewal in 8268 seconds.