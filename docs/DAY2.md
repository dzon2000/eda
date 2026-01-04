# Event contracts & schema evolution (Avro + Schema Registry)

## Goal

- Define explicit domain events
- Enforce compatibility rules
- Prove that breaking changes are blocked
- Understand how events evolve safely

## Best practices for events

- Focus on **facts**, not commands. eg. OrderCreated, OrderPaid

## Create a schema registry

```shell
curl -X POST \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data '{
    "schema": "'"$(jq -c . schemas/order-created.avsc)"'"
  }' \
  http://localhost:8081/subjects/orders.v1-value/versions
```

Unfortunately this does not work fully, as the API expects the escaped version of a JSON in string. TODO: try to find the solution.

[API Reference](https://docs.confluent.io/platform/current/schema-registry/develop/api.html#subjects)

```shell
curl -X POST   -H "Content-Type: application/vnd.schemaregistry.v1+json"   --data '{"schema": "'"{\\\"type\\\":\\\"record\\\",\\\"name\\\":\\\"OrderCreated\\\",\\\"namespace\\\":\\\"io.pw.orders.v1\\\",\\\"fields\\\":[{\\\"name\\\":\\\"orderId\\\",\\\"type\\\":\\\"string\\\"},{\\\"name\\\":\\\"customerId\\\",\\\"type\\\":\\\"string\\\"},{\\\"name\\\":\\\"amount\\\",\\\"type\\\":\\\"double\\\"},{\\\"name\\\":\\\"createdAt\\\",\\\"type\\\":\\\"string\\\"}]}"'"}' http://localhost:8081/subjects/orders.v1-value/versions
```

`--trace-ascii /dev/stdout`

Once created we can verify things with: 

```shell
curl http://localhost:8081/subjects
curl http://localhost:8081/subjects/orders.v1-value/versions
```


## Compatibility rules

[Docs](https://docs.confluent.io/platform/current/schema-registry/fundamentals/schema-evolution.html)

There are multiple types of compatibility rules:

- BACKWARD (default) - can go back one generation
- BACKWARD_TRANSITIVE - can go all the way back to first generation
- FORWARD
- FULL

Backward compatibility means that consumers can proceed data written by producers of old generations. Forward means that data produced by new publishers can be proceed with the old consumers.

```shell
curl -X PUT \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data '{"compatibility":"BACKWARD"}' \
  http://localhost:8081/config/orders.v1-value
```

Why BACKWARD:
- New consumers can read old events
- Allows safe evolution

## Go development

First we will ire Go producer with Schema Registry.

Dependencies:

```shell
go get github.com/segmentio/kafka-go
go get github.com/linkedin/goavro/v2
```

# Day 1 issues resolved

So it seems like the group creation was broken because the auto creation of kafka internal topics (`__consumer_offsets`) was not successful. Why? because apparently my cluster settings was wrong for such a simple setup (1 node). 

## ChatGPT explanation:

Kafka tries to create __consumer_offsets with:

- replication factor
- partitions

If those are **invalid** for your cluster, creation fails.

Common mistake in single-broker setups

Default:
```shell
offsets.topic.replication.factor=3
```

But you have:
- 1 broker

Kafka cannot create the topic → consumer groups break.

### Resolution

Use below settings for replication:

```shell
KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
KAFKA_TRANSACTION_STATE_LOG_MIN_ISR: 1
```

With that, the group creation and internal topics works just fine.

## Links

- [Kafka Consumer Groups and Offsets](https://axual.com/blog/kafka-consumer-groups-and-offsets-what-you-need-to-know)
- [Kafka Auto Offset](https://quix.io/blog/kafka-auto-offset-reset-use-cases-and-pitfalls)
