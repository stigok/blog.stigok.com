---
layout: post
title: "Combining bash scripts and Node.js process streams"
date: 2017-12-02 17:46:34 +0100
categories: bash nodejs
redirect_from:
  - /post/combining-bash-scripts-and-nodejs-process-streams
---

The reason I'm mixing these two is because I like the
synchronous nature of shell scripts, the JavaScript
syntax and Node.js process handling and API in general.

I want to show a way to control multiple processes and
handling their combined output in a very simple manner.

Here I have a bash shell script that prints a message
indefinitely with a fixed delay in between messages.

    # echo-delay.sh

    #!/bin/bash
    # Print <msg> to stdout every <delay> seconds
    # Arguments: <msg> <delay>
    set -eux

    while true; do
      echo "$1"
      sleep $2 
    done

My Node.js script will spawn *x* amount of these in order to
test and demonstrate merging of output streams.

    # index.js

    const execFile = require('child_process').execFile
    const es = require('event-stream')
    
    const script = './echo-delay.sh'
    const arr = [1, 2, 3, 4, 5]
    
    const streams = arr
      .map((num) => execFile(script, [num, num]))
      .map((proc) => proc.stdout)
    
    es.merge(streams)
      .pipe(process.stdout)

I'm using `event-stream` to merge an array of `ReadStreams`
into a single stream. If I was to watch heavier processes,
I could create a `Cluster` which enables me to have a process
pool spanned across multiple physical CPU cores.

Running the above Node script will give me the output of all
the spawned processes:

    $ node index.js
    1
    2
    3
    4
    5
    1
    1
    2
    1
    3
    1
    2
    4

To make the output a bit visually appealing, lets mix it up with
bars instead of numbers.

    # bars.js

    const execFile = require('child_process').execFile
    const es = require('event-stream')

    const numProcesses = Number(process.argv[2])
    const script = './echo-delay.sh'
    const repeat = (str, n) => Array(n).fill(str).join('')

    const streams = Array(numProcesses).fill('-')
      .map((c, i) => execFile(script, [repeat(c, i + 1), i + 1]))
      .map((proc) => proc.stdout)

    es.merge(streams)
      .pipe(process.stdout)

Let's run the above script, spawning 10 processes

    $ node bars.js 10
    -
    --
    ---
    ----
    -----
    ------
    -------
    --------
    ---------
    ----------
    -
    -
    --
    -
    ---
    -
    --
    ----
    -
    -----
    -
    ---
    --
    ------
    -
    -------
    -
    --
    ----
    --------
    ---
    -
    ---------

The nice thing about this is that I can have a single
script handle multiple processes with a fairly small amount
of code. Now I can add timestamps to every line of output
by using `es.split()` to split on all newlines, then  `es.through`
to easily create a `DuplexStream` function which acts on all
lines printed to stdout.

By editing the last lines of code in `index.js`, I can add a timestamp
to all lines of output pretty easily.

    # timestamps.js

    const execFile = require('child_process').execFile
    const es = require('event-stream')

    const script = './echo-delay.sh'
    const arr = [1, 2, 3, 4, 5]

    const streams = arr
      .map((num) => execFile(script, [num, num]))
      .map((proc) => proc.stdout)

    es.merge(streams)
      .pipe(es.split())
      .pipe(es.map((line, cb) => {
        const now = new Date()
        cb(null, `${now.toLocaleTimeString()}: ${line}\n`)
      }))
      .pipe(process.stdout)

Running this script it's easier to see the timing of the functions
without looking at it live:

    $ node timestamps.js 
    7:39:18 PM: 1
    7:39:18 PM: 2
    7:39:18 PM: 3
    7:39:18 PM: 4
    7:39:18 PM: 5
    7:39:19 PM: 1
    7:39:20 PM: 1
    7:39:20 PM: 2
    7:39:21 PM: 1
    7:39:21 PM: 3
    7:39:22 PM: 1
    7:39:22 PM: 2
    7:39:22 PM: 4
    7:39:23 PM: 1
    7:39:23 PM: 5
    7:39:24 PM: 1
    7:39:24 PM: 3
    7:39:24 PM: 2


## References
- https://github.com/dominictarr/event-stream
- https://nodejs.org/dist/latest-v6.x/docs/api/child_process.html#child_process_child_process_execfile_file_args_options_callback