---
layout: post
title: "Generating a pseudorandom password or string in Linux bash"
date: 2017-02-26 00:12:17 +0100
categories: linux shell
redirect_from:
  - /post/generating-a-pseudorandom-password-or-string-in-linux-bash
---

Define a function in e.g. `~/.bashrc`

    genpasswd() {
      tr -dc A-Za-z0-9 < /dev/urandom | head -c ${1:-36} | xargs
    }

Where 36 is default length if no parameter is given

## Usage

    $ genpasswd
    GVQ3ZHqrBRDzB1QwASA9uk6YsZPto2GWeRWR

    $ genpasswd 7
    qvPWx7N

## References

- http://www.shellhacks.com/en/Generating-Random-Passwords-in-the-Linux-Command-Line