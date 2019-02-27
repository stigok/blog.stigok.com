---
layout: post
title: "kafka-python no broker available with kafka on localhost"
date: 2017-12-25 02:01:36 +0100
categories: kafka python docker
redirect_from:
  - /post/kafkapython-no-broker-available-with-kafka-on-localhost
---

All of a sudden I was having problems getting a script based on `kafka-python` to work properly. Even the example code from the upstrem repo didn't work. However, I still had success with my Node.js scripts, so the problem had to be with the Python lib.

Immediately when creating a new `KafkaConsumer` or `KafkaProducer` an exception was thrown stating

    No brokers available

Both the consumer and producer constructors in `kafka-python` accepts a [long list of keyword arguments](https://kafka-python.readthedocs.io/en/master/apidoc/KafkaProducer.html#kafka.KafkaProducer), one of which is `api_version`. Apparently, when no `api_version` is set, an attempt is made to automatically detect the kafka version of the broker.

The stack trace printed along with the exception also hints about where the exception originated from, and it is when attempting to detect a version.

Specify the version explicitly, and keep on programming

    KAFKA_HOSTS = ['localhost:9092']
    KAFKA_VERSION = (0, 10)
    producer = KafkaProducer(bootstrap_servers=KAFKA_HOSTS, api_version=KAFKA_VERSION)

    topic = "foobars"
    consumer = KafkaConsumer(topic, bootstrap_servers=KAFKA_HOSTS, api_version=KAFKA_VERSION)

This may be a problem between `kafka-python` and the current Kafka installation I am using within a docker image, [wurstmeister/kafka-docker](https://github.com/wurstmeister/kafka-docker/), but I have not bothered to confirm this at this time. At least it works!

## References

- https://stackoverflow.com/a/40282989/90674