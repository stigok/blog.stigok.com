---
layout: post
title: "nginx set location of /favicon.ico"
date: 2017-04-19 15:57:55 +0200
categories: nginx html
redirect_from:
  - /post/nginx-set-location-of-faviconico
---

Browsers try to get `/favicon.ico` by default if nothing else is specified with the `<link rel="shortcut icon" href=...` tag. So because I like to keep the root clean, I put the favicon in the /img folder.

    location = /favicon.ico {
        alias /var/www/mysite/img/favicon.ico;
    }

This is kinda nice, because I don't have to do anything with the HTML.