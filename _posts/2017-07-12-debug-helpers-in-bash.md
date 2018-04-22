---
layout: post
title: "Debug helpers in bash"
date: 2017-07-12 16:47:35 +0200
categories: bash debugging
redirect_from:
  - /post/debug-helpers-in-bash
---

`set -x` prints all variables you set, along with all commands you run. Invaluable while debugging my bash scripts.

`set -e` will stop the script on any errors.

`man set` will give you the rest.