---
layout: post
title: "How to test passphrase of openssl rsa certificate keyfile"
date: 2017-06-22 18:08:44 +0200
categories: openssl bruteforce passphrase bash
redirect_from:
  - /post/how-to-test-passphrase-of-openssl-rsa-certificate-keyfile
---

I needed a way to quickly test a lot of different passphrases to a passphrase-protected certificate.key file. So I started out with the slow approach

    $ openssl rsa -check -in certificate.key
    Enter pass phrase for certificate.key:

Okay, now I want to test them in fast succession. The below snippet will ask you for the password until `openssl` exits with 0

    $ $(exit 1); while [ $? -ne 0 ]; do openssl rsa -check -in certificate.key; done

At least now I can just keep on typing until I find the right one.

The first `$(exit 1)` is to make the initial `$?` check return something else than zero. Maybe bash has a do...while too, but this works.