---
layout: post
title:  "Building my Jekyll blog in Docker using Ubuntu"
date:   2020-08-31 00:02:02 +0200
categories: jekyll docker
excerpt: I need Python to build my Jekyll site and the official Docker image to build Jekyll sites looks like a mess to me.
#proccessors: pymd
---

## Preface

I need Python to build my Jekyll site and the official Docker image to build Jekyll sites looks like a mess to me.

## Building Jekyll in Docker

The official image uses `alpine` as base, but I eventually found out I might
get less of a headache using `ubuntu:20.04` instead. That gives me good support
for both Ruby and Python 3.8 without too much hassle.

I'm not going to use the Jekyll HTTP server, but instead copy the site into
a slim `nginx` image using Docker multi-stage builds.

```dockerfile
FROM ubuntu:20.04 AS base
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
               ruby-full build-essential zlib1g-dev \
               git python3 \
    && rm -rf /var/lib/apt/lists/*
    && groupadd -g 1000 jekyll \
    && useradd -mu 1000 -g jekyll jekyll

USER jekyll
WORKDIR /home/jekyll
ENV GEM_HOME=/home/jekyll/gems \
    PATH="/home/jekyll/gems/bin:${PATH}"
RUN gem install jekyll bundler

ADD Gemfile Gemfile.lock ./
RUN bundle install

ADD . .
RUN bundle exec jekyll build -d _site -t

FROM nginx:1.19-alpine
COPY --from=base /home/jekyll/_site /var/www/html
```

This gives for a very small image in the end.

Take a look at the [current source of this Dockerfile][1] for updates.
It's the very file that builds this blog.

## References
- https://github.com/docker-library/mongo/blob/master/4.4/Dockerfile
- https://jekyllrb.com/docs/installation/ubuntu/
- https://hub.docker.com/r/jekyll/builder/tags
- https://github.com/stigok/jekyll-docker/blob/master/repos/jekyll/Dockerfile
