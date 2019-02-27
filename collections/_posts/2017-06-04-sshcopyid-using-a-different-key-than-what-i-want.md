---
layout: post
title: "ssh-copy-id using a different key than what I want"
date: 2017-06-04 21:49:23 +0200
categories: ssh ssh-copy-id pki
redirect_from:
  - /post/sshcopyid-using-a-different-key-than-what-i-want
---

I have created some different keys for different purposes, but ssh-copy-id seems to always want to use the last one I created. It picks the one with the most recent timestamp, so `touch`ing the one you want to use by default will solve it.

    $: touch ~/.ssh/id_rsa{,.pub}