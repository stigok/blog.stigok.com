---
layout: post
title: "passwordstore passdemenu and type automatically in i3wm"
date: 2017-12-14 15:03:36 +0100
categories: i3 passwordstore security
redirect_from:
  - /post/passwordstore-passdemenu-and-type-automatically-in-i3wm
---

I start `passdmenu` with a key combination handled by `i3` window manager for storing my passwords securely. On Arch Linux, this package can be installed with `pacman -S passdmenu`.

In order to make typing work, I need `xdotools`

    # pacman -S xdotools

To make `passdmenu` appear on a specific keyboard combination, update the i3 config:

    # ~/.config/i3/config (redacted version)
    set $mod Mod4
    bindsym $mod+y exec passdmenu --store $PASSWORD_STORE_DIR_PRIVATE --type

I am using `--type` to make it type passwords automatically. If you have not installed `xdotools`, it will silently fail to show, with to apparent errors.

General rule of thumb when binding things in i3 is to
- Make sure the command works from your terminal
- Make sure the full system environment has been loaded **before** i3 is started, or else you will not have the same set of system variables available