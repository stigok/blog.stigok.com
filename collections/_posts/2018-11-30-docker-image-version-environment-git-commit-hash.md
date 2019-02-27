---
layout: post
title:  "Read image version from environment or file inside Docker image in a running container"
date:   2018-11-30 16:06:01 +0100
categories: docker ci
---

## Introduction

We wanted to be able to know the current image version from within the running
container to version our socket.io connections.

We are tagging all our automatically built images with the current git commit
hash, and now we want to either read it from environment or read it from a
file from within the docker image itself.

## Dockerfile

Within the Dockerfile, describing your image, there are [lots of options][] to
control how your image is being built.

- You can set [build arguments][] (`ARG`) that are expected to be
  passed to the Docker daemon during build time.

- Values to [environment variables][] (`ENV`) can be set during build time and can
  inherit the values of build arguments.

For our simple test, we can construct a Dockerfile like this:

```Dockerfile
FROM busybox

ARG IMAGE_VERSION=''
ENV IMAGE_VERSION ${IMAGE_VERSION}

RUN echo ${IMAGE_VERSION} >> /.version

CMD ["cat", "/.version"]
```

Let's build this image with an explicit build arg and run it

```shell
$ docker build --tag image-version-test --build-arg=1.42 .
$ docker run image-version-test
1.42
```

It successfully returns the value of the build time argument, read from the
*/.version* file from the file system of the container. Let's see what the
environment variables look like.

```shell
$ docker run --rm image-version-test printenv
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=5d5f077af6db
IMAGE_VERSION=1.42
HOME=/root
```

A danger to be aware of here is that it's always possible to override the value
of the environment variable at run time, but the contents of */.version* will
always be the same for each individual image version.


## Docker compose

When building with Docker Compose, build time arguments can be specified using
`args`. Arguments with no value specified will be replaced with values from
an environment variable with the same name.

```yaml
version: '2'

services:
  image-version-test:
    build:
      context: .
      args:
        - IMAGE_VERSION
```

I can now build the image using `docker-compose` explicitly specifying a value
for `IMAGE_VERSION` from environment.

```shell
$ export IMAGE_VERSION=2.21
$ docker-compose build image-version-test
2.21
```

And since this was a build time argument, changing environment now won't make
a difference.

```shell
$ export IMAGE_VERSION=7
$ docker-compose run image-version-test
2.21
```

## Tips

If you want to build using the current git commit hash, it can be set using

```shell
$ export IMAGE_VERSION=$(git log -1 --format=%H)
```

## References
- https://docs.docker.com/compose/compose-file/#build
- https://docs.docker.com/engine/reference/builder/#arg

[lots of options]: https://docs.docker.com/engine/reference/builder/
[build arguments]: https://docs.docker.com/engine/reference/builder/#arg
[environment variables]: https://docs.docker.com/engine/reference/builder/#env

