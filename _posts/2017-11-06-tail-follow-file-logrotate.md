---
layout: post
title: "tail follow file logrotate"
date: 2017-11-06 13:18:42 +0100
categories: tail logrotate linux
redirect_from:
  - /post/tail-follow-file-logrotate
---

On one of my servers I was running `tail -f /var/log/nginx/error.log`. But output would stop every so often since the file is rotated with logrotate, giving that path a new inode number.

So when I switched to `tail -F`, I no longer have to restart the `tail` process to make it follow the output again.

## Documentation excerpt

    -f --follow
        output appended data as the file grows
    -F
        same as --follow=name --retry
    --retry
        keep trying to open a file if it is inaccessible

## References
- `man tail`
- https://unix.stackexchange.com/a/22699/28043