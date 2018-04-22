---
layout: post
title: "Setting up fail2ban on Ubuntu Server 16.04 Xenial"
date: 2017-03-15 19:12:15 +0100
categories: ubuntu fail2ban linux draft
redirect_from:
  - /post/setting-up-fail2ban-on-ubuntu-server-1604-xenial
---

# Install

    # apt-get update
    # apt-get install fail2ban

# Configure

This is a working example configuration. All code blocks below the headers belongs the the same file.

### `/etc/fail2ban/jail.d/jail.local`

Never ban specified IP in any jails

    [DEFAULT]
    ignoreip = <ip>

If `sshd` is listening on a non-default port, specify it with the `port` option.

    [sshd]
    enabled = true
    port = 2222
    bantime = 7200 ; two hours

Ban clients that matches `filter` in the specified `logpath`

    [ufw-block]
    enabled = true
    logpath = /var/log/ufw.log
    filter = ufw-block
    findtime = 180 ; three minutes
    maxretry = 5
    bantime = 3600 ; one hours

The `recidive` jail will ban clients that have been banned before. See `/etc/fail2ban/jail.conf`

    [recidive]
    enabled = true
    bantime = 86400 ; one day


### `/etc/fail2ban/filter.d/ufw-block.local`

Match BLOCK events from `ufw`.

    [Definition]
    failregex = UFW BLOCK.* SRC=<HOST>

# Helpful commands

Get information about a jail, including banned IPs and/or hostnames

    # fail2ban-client status <JAIL>

Unban an IP address

    # fail2ban set <JAIL> unbanip <IP>