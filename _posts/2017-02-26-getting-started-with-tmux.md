---
layout: post
title: "Getting started with tmux"
date: 2017-02-26 00:18:43 +0100
categories: linux shell
redirect_from:
  - /post/getting-started-with-tmux
---

# Getting started with tmux

## Install

Install with your package manager of choice. For Arch, it's `pacman`

    # pacman -S tmux

... but if you're on Ubuntu, you would typically use `apt-get` instead

    # apt-get install tmux

## Start

    $ tmux

## Detach from session

CTRL+B, D

## View keyboard shortcuts

    CTRL+B, ?

`bind-key` is CTRL+B by default