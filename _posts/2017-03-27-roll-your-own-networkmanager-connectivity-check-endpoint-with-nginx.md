---
layout: post
title: "Roll your own NetworkManager connectivity check endpoint with nginx"
date: 2017-03-27 10:22:05 +0200
categories: nginx networkmanager
redirect_from:
  - /post/roll-your-own-networkmanager-connectivity-check-endpoint-with-nginx
---

Create a file at `/etc/NetworkManager/conf.d/11-connectivity-check.conf` that describes the endpoint

    [connectivity]
    uri=http://example.com/nm-ping
    interval=120
    #response="NetworkManager is online"

The `response` option can make NetworkManager confirm connectivity by matching the response text with this string.
But I am setting up the endpoint with a `X-NetworkManager-Status` response header instead.
Set up a site on a remote server and use a site configuration like this one:

    server {
      listen 80;
      listen [::]:80;

      server_name example.com;

      location /nm-ping {
        add_header X-NetworkManager-Status online;
        add_header Content-Type text/plain;
        return 200 'NetworkManager is online';
      }
    }

Using HTTP since NetworkManager complains when it's configured with HTTPS:

    NetworkManager[467]: <warn> connectivity: use of HTTPS for connectivity checking is not reliable and is discouraged (URI: https://example.com/nm-ping)