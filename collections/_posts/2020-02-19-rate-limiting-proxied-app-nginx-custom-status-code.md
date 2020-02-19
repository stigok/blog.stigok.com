---
layout: post
title:  "Rate limiting requests to an app that is (reverse) proxied behind nginx"
date:   2020-02-19 14:18:16 +0100
categories: nginx
---

I'm exposing an API endpoint to the public internet where my program is
consuming another provider's API. I don't want to get banned for hitting
it too much, so I want to rate limit incoming request on a per-IP basis.

Using nginx, this was surprisingly easy. My server is a CentOS 7 with nginx/1.16.1
installed from the official package repositories.

I am configuring [my public transport timetable app][ruterstop] to be served under
a sub-directory on my server.

The configuration looks like this:

```nginx
limit_req_zone $binary_remote_addr zone=ruterstop:10m rate=1r/s;

server {
  listen 80;

  [redacted]

  location / {
    [redacted]
  }

  location /ruterstop/ {
    limit_req zone=ruterstop;
    limit_req_status 429;
    proxy_pass http://127.0.0.1:4000/;
    proxy_http_version 1.1;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  }
}
```

The first line above defines a zone that can be re-used within multiple `location`s.
The arguments to `limit_req_zone` should be defined as follows ([reference](http://nginx.org/en/docs/http/ngx_http_limit_req_module.html#variables)):

```nginx
limit_req_zone key zone=name:size rate=rate [sync];
```

I want to limit requests based on the source IP address ([`$binary_remote_address`](https://nginx.org/en/docs/http/ngx_http_core_module.html#var_binary_remote_addr)),
save that address in a 10MB large database file (`10m`) and set the limit per client to 1 request per second (`1r/s`):

```nginx
limit_req_zone $binary_remote_addr zone=ruterstop:10m rate=1r/s;
```

Then you can reference that zone as part of a `location` directive:

```nginx
location /a-rate-limited-path/ {
  limit_req zone=ruterstop;
  return 200 "You didn't hit the limit yet!";
}
```

Now, this works great. But the clients passing the limit would receive a `503 Service Unavailable`.
This is probably fine for most, but I'd like them to receive a `429 Too Many Requests` instead just
to make it clear that they are rate limited. This is where `limit_req_status` comes in:

```nginx
location /a-rate-limited-path/ {
  limit_req zone=ruterstop;
  limit_req_status 429;
  return 200 "You didn't hit the limit yet!";
}
```

[MDN also states](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) that a `429` could
be accompanied by a `Retry-After` header to tell the client when she's within the limit range again.
I think this should be possible by using a combination of [`error_page`][error_page] and
[`@named_location`][location] in nginx, but I'll let that remain as an excersise for those who need it.

The working (somewhat redacted) configuration was posted at the beginning of this post.

## References

[error_page]: https://nginx.org/en/docs/http/ngx_http_core_module.html#error_page
[location]: https://nginx.org/en/docs/http/ngx_http_core_module.html#location
[ruterstop]: https://github.com/stigok/ruterstop

- <http://nginx.org/en/docs/http/ngx_http_limit_req_module.html#variables>
- <https://nginx.org/en/docs/http/ngx_http_core_module.html#var_binary_remote_addr>
- <https://gist.github.com/ipmb/472da2a9071dd87e24d3>
