---
layout: post
title: "Run a command every n seconds in bash shell"
date: 2017-02-26 00:20:35 +0100
categories: linux shell utilities
redirect_from:
  - /post/run-a-command-every-n-seconds-in-bash-shell
---

> This post was later superseded by `watch`, and `while true; do $thing; sleep 1; done`

# Bash run a command every n seconds

    #!/bin/bash
    # Save this file at e.g. /usr/bin/shevery

    INTERVAL=${INTERVAL:-1}

    if [ -z "$1" ]; then
      >&2 echo "Missing command to run"
      exit 1
    fi

    while true; do
      # Run command
      $@
      sleep $INTERVAL
    done

## Usage

    $ shevery date
    Fri Feb 24 10:25:45 CET 2017
    Fri Feb 24 10:25:46 CET 2017
    Fri Feb 24 10:25:47 CET 2017
    Fri Feb 24 10:25:48 CET 2017
    Fri Feb 24 10:25:49 CET 2017
