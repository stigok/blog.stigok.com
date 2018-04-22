---
layout: post
title: "Unable to log into ssh after useradd new user"
date: 2017-06-13 02:50:54 +0200
categories: ssh
redirect_from:
  - /post/unable-to-log-into-ssh-after-useradd-new-user
---

Check the journal on the target machine closely for sshd errors

    root@server# useradd -m foob
    root@server# journalctl -fu sshd
    sshd[7693]: Invalid user foob from x.x.x.x
    sshd[7673]: User foob not allowed because account is locked
    sshd[7673]: input_userauth_request: invalid user newuser [preauth]

    foob@client$ ssh foob@server

In my case, the account is locked because the user has no password. Create a new password for the user and try again.

    root@server# passwd foob

    foob@client$ ssh foob@server