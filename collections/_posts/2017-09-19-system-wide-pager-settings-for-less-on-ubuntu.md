---
layout: post
title: "system wide pager settings for less on ubuntu"
date: 2017-09-19 15:37:03 +0200
categories: ubuntu linux man less pager
redirect_from:
  - /post/system-wide-pager-settings-for-less-on-ubuntu
---

I wanted system wide default pager settings on my Ubuntu Server 16.04. This makes manual usage of `less`, the pager in `man` and all other programs that use the `PAGER` environment variable for its pager.

Update the system environment in `/etc/environment` by appending the following values:

    LESS="-iRFX"
    PAGER="less"
    MANPAGER="less"

The `-i` option makes less ignore the case of the search string, just how I want it:

> Causes searches to ignore case; that is, uppercase and lowercase are considered identical. This option is ignored if any uppercase letters appear in the search pattern; in other words, if a pattern contains uppercase letters, then that search does not ignore case.

The `-R` option causes "raw" control characters containing ANSI color escape sequences to be output in its raw form, e.g. making colors appear. This is very helpful when using git, to avoid git outputting inline raw control chars like in the image below:

![invalid interpretation of control chars in GIT_PAGER](https://public.stigok.com/img/1507118783243146973.png)

The `-F` causes less to automatically exit if the entire file can be displayed on the first screen.

The `-X` prevents less from clearing the screen when exiting, which is important when using the `-F` flag and output fits a single screen.

The `MANPAGER` variable specifies what pager `man` should use. Actually, `man` uses `PAGER` if `MANPAGER` isn't set, so this may be superflous. In this case it uses the same settings as normal less, where the options of `less` are specified at a single place.

## References
- https://askubuntu.com/a/26411/20164
- https://askubuntu.com/q/623992/20164
- `man less`