---
layout: post
title: "Split text on whitespace in terminal output"
date: 2017-11-28 15:39:52 +0100
categories: grep linux
redirect_from:
  - /post/split-text-on-whitespace-in-terminal-output
---

To split text on whitespace you can use `grep`. There's an infinite amount of ways to do this. This is one of them.

    $ echo 'string --with ###ALLKINDS### 0f ::outputs' | grep -oP '[^\s]+'
    string
    --with
    ###ALLKINDS###
    0f
    ::outputs

## References
- `man grep`