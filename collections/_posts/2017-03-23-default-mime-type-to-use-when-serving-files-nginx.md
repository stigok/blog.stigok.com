---
layout: post
title: "Default MIME type to use when serving files nginx"
date: 2017-03-23 15:26:30 +0100
categories: nginx
redirect_from:
  - /post/default-mime-type-to-use-when-serving-files-nginx
---

I am serving some shell scripts without a file extension on my nginx web server which by default was serving these files with `Content-Type: application/octet-stream`. By using the option `default_type` I can choose what MIME type to use for files that are not already specified with the `types {}` directive.

Configuring my server to serve unknown file types as `text/plain` within a specific site:

    server {
      listen 80;
      listen [::]:80;
      listen 443 ssl;
      listen [::]:443 ssl;
      root /var/www/example.com/public;
      server_name example.com;

      default_type text/plain;
    }

## Resources

- <http://nginx.org/en/docs/http/ngx_http_core_module.html#default_type>