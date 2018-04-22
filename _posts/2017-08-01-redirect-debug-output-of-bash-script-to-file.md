---
layout: post
title: "Redirect debug output of bash script to file"
date: 2017-08-01 13:52:50 +0200
categories: bash debugging linux thunar
redirect_from:
  - /post/redirect-debug-output-of-bash-script-to-file
---

I needed a way to debug my `thunar` (file manager) custom actions, as it doesn't have a way of debugging these scripts itself.
I make my thunar *Custom actions* run a script that contains the following snippet at the top:

    #!/bin/bash
    logfile=~/tmp/bash-$$.log
    exec > $logfile 2>&1
    set -x

Where all debug output of `set -x` will go into `$logfile`. Great for making sure the supplied command arguments are sent from thunar as expected, which is why most of my custom actions didn't work before.

## References

- https://stackoverflow.com/questions/11229385/redirect-all-output-in-a-bash-script-when-using-set-x