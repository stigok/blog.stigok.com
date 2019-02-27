---
layout: post
title: "Substring extraction in bash"
date: 2017-07-27 19:32:48 +0200
categories: bash
redirect_from:
  - /post/substring-extraction-in-bash
---

To extract substrings from variables in bash, we can use the following syntax

    ${variable_name:offset:length}

Caveats:
- If offset or length is negative, count from the end
- Negative values must be separated by a space away from the colon (`:`) to avoid being interpreted by bash as *default value substitution*. (See below for examples)


## Get the file extension of a filename

    $ filename=foobar.txt
    $ echo ${filename: -4}
    .txt

## Get the file name without extension

    $ filename=foobar.txt
    $ echo ${filename:0: -4}
    foobar

## References
- `man bash`