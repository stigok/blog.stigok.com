---
layout: post
title: "Start zookeeper listening on localhost only"
date: 2018-01-03 16:03:35 +0100
categories: kafka zookeeper
redirect_from:
  - /post/start-zookeeper-listening-on-localhost-only
---

I'm setting up a local kafka server with a backing zookeeper instance, but it listens on all interfaces (`0.0.0.0`) by default.

To make zookeeper listen on localhost only, use `clientPortAddress` in the `config/zookeeper.properties`. My full config now looks like this:

    dataDir=/tmp/zookeeper
    clientPort=2181
    clientPortAddress=localhost
    maxClientCnxns=0

Now, when starting zookeeper, it prints the address it's listening to

    # cd kafka_2.12-1.0.0
    # bin/zookeeper-server-start.sh config/zookeeper.properties
    [ <redacted> ]
    [2018-01-03 15:58:54,629] INFO binding to port localhost/127.0.0.1:2181 (org.apache.zookeeper.server.NIOServerCnxnFactory)