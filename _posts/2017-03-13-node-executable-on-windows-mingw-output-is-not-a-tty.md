---
layout: post
title: "Node executable on Windows mingw 'output is not a tty'"
date: 2017-03-13 16:07:38 +0100
categories: node.js windows mingw
redirect_from:
  - /post/node-executable-on-windows-mingw-output-is-not-a-tty
---

On my Windows computer I am using mingw (Git Bash) for my terminal emulator as it's the closest I've gotten to a native Linux terminal.

I am generating some HTML with a node program I've written and I'm trying to redirect stdout to a file

    $ node generate-html.js > output.html
    output is not a tty

    $ echo $?
    1

I don't know why my terminal isn't recognized as a TTY, but apparently it has something to do about how the terminal is emulated on Windows. A solution is to execute it as a command argument to bash instead.

    $ bash -c "node generate-html.js > output.html"

    $ head -n 1 output.html
    <DOCTYPE html>

Looks good :)