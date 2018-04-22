---
layout: post
title: "kafka-python verify that a topic exists"
date: 2018-02-01 15:48:52 +0100
categories: kafka python
redirect_from:
  - /post/kafkapython-verify-that-a-topic-exists
---

I was struggeling to assert that a topic existed on the Kafka broker before starting to poll for messages. Due to a bug in `kafka-python`, polling for messages from a `KafkaConsumer` on a topic that doesn't exist will cause an infinite loop with no delay, which in turn makes CPU usage sky rocket until the process is restarted.

Finally I found a solution to first assert that a topic exists before moving on to polling:

    client = KafkaClient(self.bootstrap_servers)
    broker_topics = client.topic_partitions
    instance_topics = ['topic-2', 'topic-3', 'topic-7']

    # Make sure all topics that are to be used actually exist. This prevents
    # the consumer going into an infinite loop and 100% CPU usage when it
    # attempts to poll from a non-exising topic.
    # TODO: This will most probably be fixed in later versions of kafka-python
    for topic in instance_topics:
        if topic and topic not in broker_topics:
            print("Topic '%s' does not exist. Exiting!" % topic)
            sys.exit(1)

I'm actually exiting the process here. I run this script in a docker container and want it to exit and restart when it can't settle in its environment.

## References
- https://stackoverflow.com/questions/30943129/check-whether-a-kafka-topic-exists-in-python#30945647
- https://kafka-python.readthedocs.io/en/master/apidoc/KafkaConsumer.html