---
layout: post
title: "What is the git.io link pointing to?"
date: 2017-03-08 18:35:26 +0100
categories: git.io weave curl dns ipv6
redirect_from:
  - /post/what-is-the-gitio-link-pointing-to
---

I was trying to download and install `weave` on an IPv6-only host. `git.io` has no AAAA records, so impossible to reach. On an IPv4 enabled host, I use curl to see where the URL is redirected to

The  `s` option keeps curl from outputting any transfer stats and the `I` only outputs the response headers. Pipe to grep which only returns lines matching the pattern.

    $ curl -sI https://git.io/weave | grep 'Location:'
    Location: https://github.com/weaveworks/weave/releases/download/latest_release/weave?

Bingo! It's pointing to [https://github.com/weaveworks/weave/releases/download/latest_release/weave?](), which is kind of sad, since github.com isn't IPv6 enabled either...