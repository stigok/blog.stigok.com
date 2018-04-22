---
layout: post
title: "Get numeric (octal) permissions of a file or directory"
date: 2017-11-28 02:55:43 +0100
categories: linux
redirect_from:
  - /post/get-numeric-octal-permissions-of-a-file-or-directory
---

Linux often shows permissions using `-rwxrwxrwx` style instead of numerical like `0777`. When I started using Linux, I only knew Linux permissions from using FTP servers, and thought the numbers were called *chmod*. Now I can share that `chmod` is a utility used to modify file and directory permissions.

To get a file's permissions in octal representation you can use `stat -c %a`:

    $: stat -c %a filename
    666

To get the bitmask style

    $: stat -c %A filename
    -rw-rw-rw-

## References
- `man stat`
- `man chmod`