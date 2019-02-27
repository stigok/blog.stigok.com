---
layout: post
title: "Blink the ethernet port of your network card linux"
date: 2017-05-15 20:10:34 +0200
categories: linux nic led network
redirect_from:
  - /post/blink-the-ethernet-port-of-your-network-card-linux
---

I had a box with four ethernet ports and wanted to know which was which

    $ ethtool --identify eth2

And my NIC keeps blinking! 

## References 
- <https://serverfault.com/a/321917>