---
layout: post
title: "Automatically redirect to HTTPS in nginx"
date: 2017-03-23 15:39:31 +0100
categories: nginx SSL
redirect_from:
  - /post/automatically-redirect-to-https-in-nginx
---

The `if` statement in the below site configuration will redirect all HTTP requests to its HTTPS equivalent.

    server {
      listen 80;
      listen [::]:80;
      listen 443 ssl;
      listen [::]:443 ssl;
      
      server_name example.com;
      root /var/www/example.com/public;

      ssl_certificate     /etc/letsencrypt/live/example.com/fullchain.pem;
      ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;

      if ($scheme = 'http') {
        return 301 https://$server_name$request_uri;
      }
    }

This can be done in a more modular approach using an `include` directive, but you probably get the gist of it.

## Resources

- <http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#return>
- <http://nginx.org/en/docs/ngx_core_module.html#include>