---
layout: post
title: "Set up your own what is my ip service with just nginx"
date: 2017-07-16 20:44:16 +0200
categories: nginx whatismyip
redirect_from:
  - /post/set-up-your-own-what-is-my-ip-service-with-just-nginx
---

## How *not* to do it with nginx

I wanted to host my own *What is my ip* service. At first I made a Node.js web server which returned the remote address of the request.

    const app = express()
    app.get('/', (req, res) => {
      res.send(`Your IP address is: ${res.locals.clientip}`)
    })

However, while reading the [*list of variables*](http://nginx.org/en/docs/varindex.html) portion of the [nginx documentation](http://nginx.org/en/docs/), it occured to me how overkill my Node.js solution was, especially since I was already proxying the requests to the Node.js HTTP server via nginx.

## The good way

The [`$remote_addr` ](http://nginx.org/en/docs/http/ngx_http_core_module.html#var_remote_addr) variable can simply be returned with the [`return` directive](http://nginx.org/en/docs/http/ngx_http_rewrite_module.html#return):

    return 200 $remote_addr\n;

Also appending a newline at the end there.
Below is the nginx configuration I use for serving [ip.stigok.com](https://ip.stigok.com/)

    server {
      listen 80;
      listen [::]:80;
      listen 443 ssl http2;
      listen [::]:443 ssl http2;

      server_name ip.stigok.com;

      ssl_certificate     /letsencrypt/live/ip.stigok.com/fullchain.pem;
      ssl_certificate_key /letsencrypt/live/ip.stigok.com/privkey.pem;

      keepalive_requests 0;

      location / {
        default_type text/plain;
        return 200 $remote_addr\n;
      }
    }

The service returns something similar to this

    $ curl -i https://ip.stigok.com
    HTTP/2 200 
    server: nginx
    date: Fri, 24 Nov 2017 01:16:22 GMT
    content-type: text/plain; charset=utf-8
    content-length: 13
    x-number: 42
    
    127.13.37.1