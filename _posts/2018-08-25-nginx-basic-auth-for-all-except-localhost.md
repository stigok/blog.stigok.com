---
layout: post
title:  "Basic authentication for everyone except localhost in nginx"
date:   2018-08-25 21:25:18 +0200
categories: nginx
---

## Introduction

I have a shop order system where some admin pages are password protected using
HTTP Basic Auth in nginx. Now I want to have a monitoring daemon accessing the
admin pages without it having to authenticate itself.

- The monitoring daemon will be running on localhost
- If the request comes from the loopback device, i.e. `127.0.0.1`, allow without
  authentication
- For all other remote addresses require valid credentials with HTTP Basic Auth

```
location /admin {
  satisfy any;

  allow 127.0.0.1;
  deny  all;

  auth_basic "r u l33t f00di3?";
  auth_basic_user_file /srv/foodshop-tesoro/.htpasswd;
}
```

The clue here is the [`satisfy`][satisfy-directive] directive, which can be
either `all` or `any`. Setting `any` in this case forces the request to either
stem from localhost, or to be authenticated using `auth_basic`.

## References
- https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-http-basic-authentication/

[satisfy-directive]: https://nginx.org/en/docs/http/ngx_http_core_module.html#satisfy

