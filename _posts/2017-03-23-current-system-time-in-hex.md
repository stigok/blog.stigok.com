---
layout: post
title: "Current system time in hex"
date: 2017-03-23 14:41:53 +0100
categories: linux bash
redirect_from:
  - /post/current-system-time-in-hex
---

Get current system time as HEX string

    #!/bin/bash

    # Get zero-padded hex string from date format
    # Usage: f <format>
    function f {
    	printf %02X $(date +$1)
    }

    h=$(f %-H)
    m=$(f %-M)
    s=$(f %-S)

    echo $h:$m:$s

Example

    $ date && timestring
    Thu Mar 23 14:34:02 CET 2017
    0E:22:02