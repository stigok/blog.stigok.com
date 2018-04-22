---
layout: post
title: "Set default application for MIME type with xdg-mime"
date: 2017-08-30 18:51:13 +0200
categories: arch linux xdg-mime mime
redirect_from:
  - /post/set-default-application-for-mime-type-with-xdgmime
---

My JSON and other plain text files were opening in Libre Office after I installed it. That is not something I appreciate. Here's how to change it to something else using `xdg-mime`.

See what is opening plain text files right now

    $ xdg-mime query default text/plain

See a list of available applications

    $ ls /usr/share/applications

Change default application

    $ xdg-mime default mousepad.desktop text/plain

## References
- https://wiki.archlinux.org/index.php/Default_applications#XDG_standard
- https://wiki.archlinux.org/index.php/Desktop_entries