---
layout: post
title: "Testing all substrings of a password for basic auth with curl"
date: 2017-07-21 19:29:19 +0200
categories: curl bash http
redirect_from:
  - /post/testing-all-substrings-of-a-password-for-basic-auth-with-curl
---

I had to test a substring of a password for a weird authentication issue with a local client

    $ pw=K8ZUCmecnQU5o84mHPCnNsAKS4EdONbsJDcV
    $ i=32; while true; do url="http://admin:${pw:0:$i}@10.3.2.30/"; curl -f -s -w '%{http_code}' $url; i=$((i - 1)); echo -e "\n$url"; read; done

So with this script I press enter until I see a HTTP 200, then use that URL for logging in via a browser.