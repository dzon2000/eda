# Idempotent producers and consumers

Idempotency in Kafka ensures that producing or consuming the same message multiple times has the same effect as doing it once (“exactly-once effect” at the application level).

Idempotent producers:
- Use a producer ID (PID) and sequence numbers per partition to prevent duplicates caused by retries.
- On broker side, messages with the same PID and sequence number are accepted only once.
- Enable via `enable.idempotence=true`; often combined with `acks=all` and appropriate retries.
- Gives per-partition exactly-once *delivery semantics* from the producer to Kafka (no lost or duplicated records in a given partition), assuming the broker is configured correctly.

Idempotent consumers:
- Cannot rely on Kafka alone; require application-level logic and durable state.
- Common patterns:
    - Track processed message IDs / offsets in an external store (e.g., DB) and skip already processed ones.
    - Use a transactional outbox/inbox pattern so each message’s side effects are applied only once.
    - Commit consumer offsets only after side effects are completed, or manage offsets in the same transaction as state changes (often with Kafka transactions and a database or Kafka Streams).
- Goal is that reprocessing (due to consumer restarts, retries, or rebalances) does not cause duplicated side effects (e.g., double-charging a customer).

## Goals

- Build application-level idempotency
    - Side note `segmentio/kafka-go` does not support idempotency
