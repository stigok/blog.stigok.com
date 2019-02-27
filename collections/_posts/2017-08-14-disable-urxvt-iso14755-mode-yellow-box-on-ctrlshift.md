---
layout: post
title: "Disable urxvt ISO14755 mode yellow box on ctrl+shift"
date: 2017-08-14 14:28:09 +0200
categories: urxvt unicode
redirect_from:
  - /post/disable-urxvt-iso14755-mode-yellow-box-on-ctrlshift
---

Whenever I press Ctrl+Shift when in the urxvt terminal, the key combination toggles the ISO14755 mode of the terminal, letting you enter unicode sequences manually. I have personally never used this feature myself and just got tired of dealing with it as it would toggle on while I was trying to trigger a keybord shortcut of a different background process (i3wm).

![urxvt ISO 4755 mode](https://public.stigok.com/img/1502713196362161752.png)
![urxvt keycap picture instert mode](https://public.stigok.com/img/1502713211095529058.png)

urxvt has an way to disable the mode completely with the `iso14755` option. So I put the following into my `~/.Xresources`:

    URxvt.iso14755: false

Then I reload the configuration

    $ xrdb -merge ~/.Xresources

New instances of urxvt should no longer toggle the mode on Ctrl+Shift.

## References
- https://bbs.archlinux.org/viewtopic.php?pid=1206991#p1206991