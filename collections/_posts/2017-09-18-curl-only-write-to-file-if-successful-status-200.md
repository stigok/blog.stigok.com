---
layout: post
title: "curl only write to file if successful status 200"
date: 2017-09-18 16:40:08 +0200
categories: curl
redirect_from:
  - /post/curl-only-write-to-file-if-successful-status-200
---

Make curl get the contents of a URL and write to file, but _only_ write to file if the response is successful:

    curl -s -S -f -o facebook-feed.json "$facebook"

- `-s` keeps curl quiet by hiding progress meter and error messages
- `-S` shows an error message if it fails (stderr)
- `-f` Fail  silently (no output at all) on server errors, keeping stdout clean
- `-o` specifies an output file