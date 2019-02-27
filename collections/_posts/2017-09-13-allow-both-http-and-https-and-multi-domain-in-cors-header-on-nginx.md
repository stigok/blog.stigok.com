---
layout: post
title: "Allow both HTTP and HTTPS and multi domain in CORS header on nginx"
date: 2017-09-13 03:51:36 +0200
categories: cors nginx
redirect_from:
  - /post/allow-both-http-and-https-and-multi-domain-in-cors-header-on-nginx
---

Wanted to allow both HTTP and HTTPS to POST to the API of an old solution I've been maintaining lately. But the `Access-Control-Allow-Origin` only allows a single domain with strict scheme. This made me look around in the nginx docs again.

Since my site is also answering on multiple domains, this is pretty neat:

    server_name example.com www.example.com;
    location /api/ {
        add_header Access-Control-Allow-Origin "${scheme}://${server_name}";
        proxy_pass http://localhost:4242/;
    }

## References
- http://nginx.org/en/docs/varindex.html
- https://serverfault.com/questions/152194/merging-variable-with-string-in-config-file