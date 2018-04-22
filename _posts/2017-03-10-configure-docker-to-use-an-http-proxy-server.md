---
layout: post
title: "Configure docker to use an HTTP proxy server"
date: 2017-03-10 11:01:41 +0100
categories: docker systemd
redirect_from:
  - /post/configure-docker-to-use-an-http-proxy-server
---

There are several ways to configure this, some of which are simpler or harder, but I think this was an interesting way to do it.

    # mkdir -p /etc/systemd/system/docker.service.d
    # echo -e '[Service]\nEnvironment="HTTP_PROXY=http://[::1]:8080/"' > /etc/systemd/system/docker.service.d/http-proxy.conf

Reload configuration and restart the docker daemon

    # systemctl daemon-reload
    # systemctl restart docker

Verify that docker is run with the new environment

    # systemctl show --property Environment docker
    Environment=GOTRACEBACK=crash DOCKER_HTTP_HOST_COMPAT=1 HTTP_PROXY=http://[::1]:8080/

## References

  - [https://docs.docker.com/engine/admin/systemd/#http-proxy]()