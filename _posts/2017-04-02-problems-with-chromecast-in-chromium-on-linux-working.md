---
layout: post
title: "Problems with Chromecast in Chromium on Linux (working!)"
date: 2017-04-02 23:37:55 +0200
categories: chromium chromecast arch avahi
redirect_from:
  - /post/problems-with-chromecast-in-chromium-on-linux-working
---

After too long without being able to cast to my speakers, it was time to get it up and running. Started off with starting the `avahi-daemon` and getting `mdns` up and running. That made me able to find the chromecast device in the local mDNS group, but Chromium saw nothing.

![Chromium Cast: no cast destinations found](https://s.42.fm/img/img-1491204315.png)

Assuming you have chromium installed already. My version information:

    $ chromium --version
    Chromium 57.0.2987.133

Install packages

    # pacman -S avahi nss-mdns

Edit `nsswitch.conf` and add `mdns_minimal [NOTFOUND=return] ` just before `resolve` on the `hosts:` line

    [...]
    hosts: files mymachines mdns_minimal [NOTFOUND=return] resolve [!UNAVAIL=return] dns myhostname
    [...]

Configuring my firewall. This is only neccessary if the default outgoing policy is `deny` or `reject`.

    # ufw allow out to any port 5353 proto udp

Start Chromium with custom flag to enable *Media Router Component Extension* or enable it in Chrome through `chrome://flags/#load-media-router-component-extension`. The command line argument may or may not work, but the flag sets it for sure.

    $ chromium --load-media-router-component-extension

This is where chromium segfaults and dumps its core

    kernel: Networking Priv[1668]: segfault at 10 ip 000055a8a812aac8 sp 00007ff4c872ed30 error 4 in chromium[55a8a6331000+8833000]

If you previously set the flag manually in `chrome://about:flags` you can revert it by editing `~/.config/chromium/Local\ State`

    $ sed -e 's/load-media-router-component-extension@0/load-media-router-component-extension@1/' ~/.config/chromium/Local\ State

This is the current state of Chromecast support in Chromium. I have seen some people around the net getting it to work, but their setup is uknown. This post will be updated when I get it to work properly again.

## Working in Chromium 59.0.3071.115 (as of 2017-07-10)

![Chromium](https://public.42.fm/1499719245183894519.png)

Works again!

## References

- <http://peter.sh/experiments/chromium-command-line-switches/>
- <https://wiki.archlinux.org/index.php/Avahi#Obtaining_IPv4LL_IP_address>
- <https://wiki.archlinux.org/index.php/Chromium#Chromecasts_in_the_network_are_not_discovered>
- <https://bugs.archlinux.org/task/51832>