# Objectives

- Make docker compose for Kafka and schema registry

## What I have learnt

- Kafka is distributed streaming for pub/sub events
- Kafka has a lot of concepts - need to dig into it more
- bitnami images are no longer free. Had to use official apache/kafka image
- Starting a base image rather simple. Look for `docker logs docker-kafka-1 | grep --color=auto -i "Kafka Server started"`
- Need to dig more into partitioning, offsets and groups
    - This is interesting combination as I had some issues reading messages with consumer
- Kafka comes with a lot of tools that help interact with the cluster
- Kafka no longer uses Zookeeper but rather implements KRaft (raft algorithm) for orchestration
>Raft offers a generic way to distribute a state machine across a cluster of computing systems, ensuring that each node in the cluster agrees upon the same series of state transitions.

## Worth to remember

- Make sure to add kafka scripts to PATH: `PATH="/opt/kafka/bin/:$PATH"`
- Describe broker: `kafka-configs.sh --bootstrap-server localhost:9092 --describe --broker 1`
    - This will not print anything if there is no dynamic config which is OK
- Create topic:
```shell
kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create \
  --topic orders.v1 \
  --partitions 3 \
  --replication-factor 1
```
- Describe topic: `kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic orders.v1`
- Run a producer script:
    - This will start interactive shell. Each enter is a new message sent to a topic
```shell
kafka-console-producer.sh \
  --bootstrap-server localhost:9092 \
  --topic orders.v1
```
- Run a consumer:
    - This will print messages sent to a topic. For some reason it didn't work until I used `--partition 0` switch
```shell
kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic orders.v1 \
  --from-beginning
```

## Day 1 outcome

- Kafka cluster and topic created
- Producer script can send messages to a topic
- Consumer get process messages*
- Docker compose up and ready

## Final thoughts

I need to have a closer look to partitions, offsets and groups. For some reason consumer didn't process any message until I used explicit `--partition 0` switch. This does not make sense to me as what if the actual parition is not 0? As far I found out it should be managed by using groups but for some reason the group was not created and thus the consumer didn't process any messages.