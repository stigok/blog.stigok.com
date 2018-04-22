---
layout: post
title: "urxvt backspace and other control keys not working on remote host in ssh"
date: 2017-05-13 18:52:03 +0200
categories: ssh urxvt rxvt-unicode terminfo
redirect_from:
  - /post/urxvt-backspace-and-other-control-keys-not-working-on-remote-host-in-ssh
---

As I was connecting via SSH to a remote host running Arch ARM the terminal was apparently interpreting all control keys as whitespace characters. As usual, the Arch Linux wiki knew exactly what was going on.

First check your current terminal

    $ echo $TERM
    rxvt-unicode-256-colors

Then copy the correct terminfo file to the remote host

    $ ssh remotehost mkdir -p .terminfo/r/
    $ scp /usr/share/terminfo/r/rxvt-unicode-256color remotehost:.terminfo/r/

Then connect again to the host and enjoy a (hopefully) working terminal again.

## References
- <https://wiki.archlinux.org/index.php/Rxvt-unicode#Remote_hosts>