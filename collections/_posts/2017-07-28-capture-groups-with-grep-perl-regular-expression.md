---
layout: post
title: "Capture groups with grep perl regular expression"
date: 2017-07-28 10:53:35 +0200
categories: grep linux arch regex
redirect_from:
  - /post/capture-groups-with-grep-perl-regular-expression
---

I got the example idea from a [Stack Overflow question](https://unix.stackexchange.com/a/192852/28043) where a method for getting the IP address of the primary DNS server was used like this

    $ cat /etc/resolv.conf | grep -i '^nameserver' | head -n1 | cut -d ' ' -f2

I the snippet reads pretty clear, but the pipeline is rather long. At first I thought I could at least drop the `head` part, as I can limit `grep` amount of matched lines with `-m`

    $ cat /etc/resolv.conf | grep -i -m 1 '^nameserver' | cut -d ' ' -f2

But I think I can do better. `grep -P` makes use of *Perl-compatible regular expression (PCRE)* patterns and opens up for use of capture groups with `\K`. To return *only* the contents of the capture group, use `-o`. And I know the `nameserver` portion is lowercase, so I don't want case-insensitive search; dropping `-i`

    $ cat /etc/resolv.conf | grep -P -o -m 1 '^nameserver \K\S+'

Where the regex pattern reads as follows:

    #  `^` start of line
    #  `nameserver ` match exact string contents
    #  `\K` start capture group
    #  `\S+` one or more non-whitespace characters

But.. I think I can do even better. Since the beginning this snippet, it has been subject to a *useless use of `cat`*. `grep` takes a *file* argument after the pattern

    $ grep -Pom 1 '^nameserver \K\S+' /etc/resolv.conf

## References
- `man grep`
- https://unix.stackexchange.com/a/192852/28043