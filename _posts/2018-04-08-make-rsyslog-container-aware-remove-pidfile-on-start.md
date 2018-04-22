---
layout: post
title: "Make rsyslog container aware remove pidfile on start"
date: 2018-04-08 20:04:32 +0200
categories: rsyslog docker
redirect_from:
  - /post/make-rsyslog-container-aware-remove-pidfile-on-start
---

There was [a commit](https://github.com/rsyslog/rsyslog/commit/d4c6c5f468ed40a0a6f614614822750d9c0255c9) to rsyslogd that made the [rsyslog daemon container-aware](https://www.rsyslog.com/doc/master/containers/container_features.html) by checking if it was running as PID 1. If so, it will avoid writing a pidfile to disk. This is very handy when using it in a container environment where you might have a lot of container restarts or recreations, and you don't want the daemon to check if a pidfile exists, because when it does, rsyslogd will refuse to start.

I am using Debian Stretch in my Docker image, but a version of rsyslog including this commit does not exists in the Debian apt repositories.

A work-around is to delete the pidfile manually whenever a container starts up by using a custom `ENTRYPOINT`, which is what I ended up with. In which case, here's a file to include in your custom rsyslog image `bin/entrypoint.sh`:

    #!/bin/sh
    # Clear the pid file from old runs. It doesn't restart otherwise.
    # NOTE: This won't be necessary after upgrading to version 8.34.0 of rsyslogd
    rm /var/run/rsyslogd.pid
    
    # Run command after configuration
    exec $@

Then reference this file in the `Dockerfile`:

    # Use stretch-backports to get rsyslog version with mongodb uristr option
    FROM debian:stretch-backports
    
    RUN apt update && apt -t stretch-backports install -y \
        rsyslog \
        rsyslog-relp \
        rsyslog-mongodb
    
    # This must be updated if more protocols are added
    EXPOSE 514/tcp
    
    # rsyslogd will die on QUIT, TERM or INT (man rsyslogd)
    STOPSIGNAL SIGTERM
    
    COPY rsyslog.conf /etc/rsyslog.conf
    
    COPY bin/configure.sh /usr/bin
    ENTRYPOINT ["configure.sh"]
    
    # This will enable debug output toggling with SIGUSR1 (requires -d)
    #ENV RSYSLOG_DEBUG=DebugOnDemand
    
    # -n : Avoid auto-backgrounding
    # -d : Enable debug mode
    CMD ["rsyslogd", "-n", "-d"]

An alternative to this is to build your own image of rsyslog from source, but I didn't want to do that.