---
layout: post
title: "ssh directly into tmux"
date: 2017-10-30 16:42:48 +0100
categories: ssh tmux
redirect_from:
  - /post/ssh-directly-into-tmux
---

I use this all the time. If there are no active sessions, it most likely means my server has rebooted unexpectedly.

    $: ssh user@myserver -t tmux attach -t 0