---
layout: post
title: "ufw allow rules on multiple interfaces"
date: 2017-03-07 17:24:57 +0100
categories: ufw iptables linux
redirect_from:
  - /post/ufw-allow-rules-on-multiple-interfaces
---

I was trying to set up some rules on multiple interfaces at once with `ufw`, but had some unsuccessful attempts. Between 

This did not open for outgoing traffic on `tun0`:

    $ sudo ufw allow out on tun0,tun1,tun2 to any
    $ echo $-
    0

But all of a sudden I was able to ping out on `tun0` with this one:

    $ sudo ufw allow out on tun0
    $ echo $-
    0

Then to see if my previous `tun0,tun1,tun2` wasn't what I thought it was:

    $ sudo ufw allow out on qwerty
    $ echo $-
    0

## What I learned about ufw

- ufw doesn't care whether an interface name exists or not
- interface names can contain commas
- it's not possible to configure multiple interfaces in a single rule