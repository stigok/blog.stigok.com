---
layout: post
title: "Vim indentation for specific file type extensions"
date: 2017-12-25 02:28:36 +0100
categories: vim
redirect_from:
  - /post/vim-indentation-for-specific-file-type-extensions
---

I am very fond of 2 spaces for indentation in most of my files, but I want to follow the convention with using a width of 4 when coding in Python.

The easiest way I could find how was to add a line in `~/.vimrc`

    autocmd FileType python setlocal shiftwidth=4

In my own setup (which may not apply to you), I've set the global settings as such:

    set shiftwidth=2
    set softtabstop=2

## References
- https://stackoverflow.com/a/159066/90674