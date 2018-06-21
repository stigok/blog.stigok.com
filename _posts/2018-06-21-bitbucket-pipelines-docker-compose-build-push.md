---
layout: post
title:  "Using docker-compose build and push in Bitbucket Pipelines"
date:   2018-06-21 08:14:23 +0200
categories: docker
---

## Introduction
A project I'm working on is divided into several different docker images.
For automatic image builds and pushing to remote container registries, we
are using Bitbucket Pipelines to define the process.

Additionally we want to tag the images with two different tags. One as
the latest `nightly`, and one with the commit hash.

The file used to describe the pipelines is *bitbucket-pipelines.yml* and
should be placed in the root of your repository.

## Build and push

If I was building and pushing a single image, I would be using plain `docker`
commands.

```
$ export COMMIT=$(git rev-parse --short HEAD)
$ docker build my-image:nightly -t my-image:$COMMIT
$ docker push my-image:nightly
$ docker push my-image:$COMMIT
```

When I have multiple images, I have to run the above thre commands for each
of them. This is where `docker-compose` comes into the picture.

The build process using `docker-compose` is vastly simplified, but this also
only allows for a single tag per image, per command.

However, since version 3 of the docker-compose file syntax I can use variable
substition to help out in the process of tagging multiple images with the same
tag:

```
version: "3"
services:
  a:
    image: my-image:${IMAGE_TAG:-nightly}
    build: my-image/
  b:
    image: my-second-image:${IMAGE_TAG:-nightly}
    build: my-second-image/
```

Whenever I'm running `docker-compose build` without setting `IMAGE_TAG` in my
environment, it will be tagged as `nightly`. This means that in order to tag
both of these images with both tags, I can run the the same command for two
different environments.

```
$ docker-compose build
$ IMAGE_TAG=$COMMIT docker-compose build
```

Since the images are already built with the first command, the second one will
use the cache and skip the build process to simply just tag them.

Now I also want to push them to a remote repository. Since the remote
container registry URL might change, I am introducing another variable into
the compose file.

```
version: "3"
services:
  a:
  image: ${IMAGE_NAME_PREFIX}my-image:${IMAGE_TAG:-nightly}
    build: my-image/
  b:
    image: ${IMAGE_NAME_PREFIX}my-second-image:${IMAGE_TAG:-nightly}
    build: my-second-image/
```

When `IMAGE_NAME_PREFIX` isn't set, it will default to a blank string. This is
very helpful as we can now build locally on the machine without cluttering
the local docker image list with lots of remote tags, while they can still
easily be tagged with a remote registry or registry username.

For local builds, we would still run `docker-compose build` as normal, but in
the pipelines config we can load an environment variable and push remotely.

An example *bitbucket-pipelines.yml* for this solution now:

```
pipelines:
  branch:
    develop:
      - step:
          name: image-build-and-push
          services:
            - docker
          caches:
            - docker
          script:
            - export
              IMAGE_NAME_PREFIX="${CONTAINER_REGISTRY}/$DOCKER_USERNAME/"
              IMAGE_TAG="${BITBUCKET_COMMIT}"
            - docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD https://${CONTAINER_REGISTRY}
            - docker-compose build && docker-compose push
            - IMAGE_TAG=nightly docker-compose build && docker-compose push
```

In the pipelines settings I will set the `CONTAINER_REGISTRY`,
`DOCKER_USERNAME` and `DOCKER_PASSWORD` environment variables and it will
build whenever I push something to my develop branch. `BITBUCKET_COMMIT` is
populated automatically and will contain the full length hash of the current
git commit.

## References
- https://docs.docker.com/compose/compose-file/#variable-substitution
- https://docs.docker.com/compose/reference/envvars/#docker_cert_path
- https://confluence.atlassian.com/bitbucket/branch-workflows-856697482.html
