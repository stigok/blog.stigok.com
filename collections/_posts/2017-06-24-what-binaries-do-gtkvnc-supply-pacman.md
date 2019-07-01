---
layout: post
title: "What binaries does a program supply? (pacman)"
date: 2017-06-24 11:27:35 +0200
categories: pacman archlinux
redirect_from:
  - /post/what-binaries-do-gtkvnc-supply-pacman
---

I installed `gtk-vnc` to connect to a locally connected Mac via VNC, but
after installation I had no idea what binaries was installed.

With `pacman`, you can use `-Ql` flags to list a package's provided files.

    $ pacman -Ql gtk-vnc | grep bin
    gtk-vnc /usr/bin/
    gtk-vnc /usr/bin/gvnccapture
    gtk-vnc /usr/bin/gvncviewer
