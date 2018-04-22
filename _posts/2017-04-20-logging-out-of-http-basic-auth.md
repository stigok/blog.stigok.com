---
layout: post
title: "Logging out of HTTP basic auth"
date: 2017-04-20 11:14:06 +0200
categories: http authentication
redirect_from:
  - /post/logging-out-of-http-basic-auth
---

If I'm already at `https://stigok.com/protected/area` which is protected with HTTP basic authentication, I can overwrite the cached credentials by sending myself to `https://anyotherstring@stigok.com/protected/area`. Note that this is a client based approach and may not work in all browsers.