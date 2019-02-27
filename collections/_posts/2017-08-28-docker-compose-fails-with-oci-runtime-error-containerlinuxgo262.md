---
layout: post
title: "Docker compose fails with oci runtime error container_linux.go:262"
date: 2017-08-28 17:51:25 +0200
categories: docker docker-compose arch
redirect_from:
  - /post/docker-compose-fails-with-oci-runtime-error-containerlinuxgo262
---

Attempting to create a pair of containers with `docker-compose` but I'm getting
error messages like

    oci runtime error: container_linux.go:262

Tool versions

    $ docker -v
    Docker version 17.06.1-ce, build 874a7374f3
    $ docker-compose -v
    docker-compose version 1.15.0, build unknown

And my Docker compose file was

    version: '3.3'

    services:
      db:
        image: postgres
        environment:
          POSTGRES_PASSWORD: p420w0rd
          POSTGRES_DB: test
        volumes:
          - postgres:/var/lib/postgresql/data

      adminer:
        image: adminer
        ports:
          - 8080:8080

    volumes:
      postgres:

But attempting to bring them up returns errors:

    $ docker-compose up --remove-orphans
    Creating postgres_adminer_1 ...
    Creating postgres_db_1 ...
    Creating postgres_adminer_1
    Creating postgres_db_1 ... error
    ERROR: for postgres_db_1  Cannot start service db: oci runtime error: container_linux.go:262: starting container process caused "process_linux.go:339: container init caused \"process_linux.go:322: running prestart hook 0 caused \\\"fork/exec /usr/bin/dockerd (deleted): no such file or directory\\\"\""
    Creating postgres_adminer_1 ... error
    ERROR: for postgres_adminer_1  Cannot start service adminer: oci runtime error: container_linux.go:262: starting container process caused "process_linux.go:339: container init caused \"process_linux.go:322: running prestart hook 0 caused \\\"fork/exec /usr/bin/dockerd (deleted): no such file or directory\\\"\""
    ERROR: for db  Cannot start service db: oci runtime error: container_linux.go:262: starting container process caused "process_linux.go:339: container init caused \"process_linux.go:322: running prestart hook 0 caused \\\"fork/exec /usr/bin/dockerd (deleted): no such file or directory\\\"\""
    ERROR: for adminer  Cannot start service adminer: oci runtime error: container_linux.go:262: starting container process caused "process_linux.go:339: container init caused \"process_linux.go:322: running prestart hook 0 caused \\\"fork/exec /usr/bin/dockerd (deleted): no such file or directory\\\"\""
    ERROR: Encountered errors while bringing up the project.

Looking around the net, I saw someone mentioning that this might come from a kernel error.
Perhaps it arised from a kernel update of which I did not reboot afterwards.

**A reboot solved the issue with no sight of the error messages.**

Looking at the package manager logs, I can see when I last updated the kernel:

    $ grep 'upgraded linux' /var/log/pacman.log
    [2017-08-21 10:22] [ALPM] upgraded linux (4.12.8-1 -> 4.12.8-2)

And my last reboot was before that

    $ journalctl --list-boots
     -1 f1e4a0f684ae4e3ca199874501e32d5f Thu 2017-08-10 12:50:08 CEST—Mon 2017-08-28 17:27:52 CEST

So that seems to actually be the case for me.