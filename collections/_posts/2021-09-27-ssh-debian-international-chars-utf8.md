---
layout: post
title:  "Internation chars (UTF-8) not displaying properly over SSH"
date:   2021-09-27 11:04:38 +0200
categories: linux debian locales
excerpt: UTF-8 chars are displayed as underscores over SSH to a Debian box
#proccessors: pymd
---

## Preface

I just installed a new Debian box and when connecting over SSH there seems to be
an issue with displaying Norwegian UTF-8 characters like `æøåÆØÅ`.

I'm pretty sure I made a mistake during my *Advanced Graphical Installation* of
Debian 11 by forgetting to configure locales.

When writing these chars (`æøå`) in a shell, they would appear as underscores in
some of my programs and in my shell, `/bin/bash`, as control characters `[^9`
and so on.
They would however appear properly while reading interactively from stdin

```
$ cat
æøå
```

I don't really know why it works through stdin (please send me a comment if you know),
but now it works again as it should everywhere with the solution below.

## Configure locales in Debian 11

So I checked my *locales* configuration and it appeared I had not installed any

```
# dpkg-reconfigure locales
```

I went ahead and selected `en_US.UTF-8 UTF-8` and `nb_NO.UTF-8 UTF-8` from the list and terminated
all my login shells, deleted the existing SSH connection and reconnected.

I'm selecting Norwegian here because I want my `LC_TIME`, `LC_PAPER` and some others
in my country's locale. You probably don't need it yourself.

Behold -- now it looks like it should again :)

## Notes

At first, I thought I only had these problems in irssi, my terminal IRC client,
but even when opening a normal `/bin/bash` shell the problem was there too.
After fixing my locales this also fixed irssi.


## References
- https://perlgeek.de/en/article/set-up-a-clean-utf8-environment
