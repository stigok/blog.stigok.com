---
layout: post
title:  "Extracting substring in bash by matching"
date:   2018-10-24 14:31:08 +0200
categories: bash
---

## Introduction

In a Bitbucket Pipelines configuration I am building a release pipeline that
matches a release branch by name `release/*`. In this build step I am building
a docker image, then in turn tagging it. I want to tag it using the version
in the branch name, just without the `release/` part.

The branch name is available in `$BITBUCKET_BRANCH`, and I can use bash to
extract a substring.

## Extract substring

Remove matching prefix pattern `${parameter##word}`

```bash
$ branch=release/branch/2.2.1; version=${branch##*/}; echo $version;
2.2.1
```

Using `##`(as above) will return the longest match.
Using a single `#` (as below) will return the shortest match.

```
$ branch=release/branch/2.2.1; version=${branch#*/}; echo $version;
branch/2.2.1
```

This means I can use the following to get the image tag I want

```
$ docker build -t my-image:${BITBUCKET_BRANCH#release/}
```

To be shorter, less explicit, alas less readable for the untrained basher

```
$ docker build -t my-image:${BITBUCKET_BRANCH#*/}
```

Readability matters the most!

## References
- http://www.tldp.org/LDP/abs/html/string-manipulation.html
- `man bash`
