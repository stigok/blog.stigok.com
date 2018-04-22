---
layout: post
title: "What binaries do gtk-vnc supply? (pacman)"
date: 2017-06-24 11:27:35 +0200
categories: pacman vnc gtk-vnc mac archlinux
redirect_from:
  - /post/what-binaries-do-gtkvnc-supply-pacman
---

Installed `gtk-vnc` to connect to a locally connected Mac via VNC. Had trouble finding out what binaries it supplied.

    $ pacman -Ql gtk-vnc | grep bin
    gtk-vnc /usr/bin/
    gtk-vnc /usr/bin/gvnccapture
    gtk-vnc /usr/bin/gvncviewer

Now I know.