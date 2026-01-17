# DLQ

Dead Letter Queue (DLQ) is a dedicated Kafka topic used to store messages
that cannot be processed successfully by consumers (e.g., due to parsing 
errors, schema mismatches, or business validation failures).

Using a DLQ allows the main processing pipeline to continue operating
without being blocked by problematic records, while preserving failed
messages for later inspection, debugging, and potential reprocessing.

## Goals

- Implement DLQ for orders
- Implement DLQ producer
- Implement mechanism to handle processing failures in original consumer

# Producer reliability

Design a producer that:
- Never loses events
- Survives crashes
- Is auditable
- Scales operationally

> Store events in the same transaction as business data

E.g. use DB to store events, and poll to send to Kafka

Example:

```sql
BEGIN;

INSERT INTO orders (...);

INSERT INTO outbox_events (
    id,
    aggregate_type,
    aggregate_id,
    event_type,
    payload,
    schema_version
) VALUES (
    gen_random_uuid(),
    'order',
    :order_id,
    'OrderCreated',
    :payload::jsonb,
    2
);

COMMIT;
```