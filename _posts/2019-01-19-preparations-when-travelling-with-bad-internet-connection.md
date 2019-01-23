---
layout: post
title:  "Travel tips for bad internet connections"
date:   2019-01-19 16:05:08 +0100
categories: lifestyle
---

I am travelling in Brazil now, in places with varying quality of internet
connectivity. Here are some tips that may make life easier the next time
it's time to travel.

1. Have as much as possible available in a terminal on a remote host with
   good internet. Accessing a remote host over SSH is both simple and
   secure, and when you end up places with poor internet, it will make
   certain tasks a lot faster.

   A shared benefit for all terminal processes over SSH is that you're only
   transferring the text content visible on the screen at any given time.

2. Have a terminal email client. Having to download all e-mails to my
   local device, both text and html bodies, before reading is not
   enjoyable to sit around waiting for. When packets get lost in
   transit, some downloads and slow connection handshakes fail.
   Having this client on the remote host with shell access will be ideal.

3. Have a local DNS caching server. Waiting for DNS lookups from a slow
   server is anoying. Set up a local server that at least caches queries,
   and possibly, also gives you stale requests when records are out of
   date to keep you from waiting unecessarily.

4. Try to keep your stuff on TCP connections. In UDP, the emitting party doesn't
   care if the data reaches the destination or not. This is a problem when
   sitting on poor connections, especially WiFi connections, where packets
   disappear in transit all the time. For example, if you can choose between
   using a VPN over UDP or TCP, go for TCP when packets get lost on reconnects
   occur frequently.

5. Download all your movies and music to local drive before you leave the
   house.

6. Be happy. Everything is going to be alright.

