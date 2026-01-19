# Consumer as stateful system

Offset ≠ State (the most important Kafka concept).

Many engineers believe:

> “If I committed the offset, the message is processed.”

This is **false*.

An offset only means:

> “The consumer read this message.”

Kafka does not know:
- if DB write succeeded
- if state was updated
- if downstream events were sent

## State is a side effects

An example of the side effect could be:
- a row inserted into PostgreSQL
- an order marked as PAID
- an email sent
- a cache updated

Kafka does not now that and does not need to know. It only tracks reading, not effects.

> Commit offsets only after state is safely persisted.

## Consumer flow

```
read message
BEGIN
  if event_id exists → skip
  apply business change
  insert event_id
COMMIT
commit offset
```

**Don't use auto commit**

There are cases where auto commit is acceptable:
- Occasional loss is acceptable
- Duplicates don’t matter
- No strong state

An example is consuming application logs or where correctness is not a key (e.g. approximation in analytics, observability etc)

## When to commit messages

There are few options:
- Per message
    - Good for business critical
    - Easy to implement/debug
    - But: high commit overhead, lower throughput
- Batch
    - Higher throughput
    - Better DB utilization (if used with state)
    - But: harder troubleshooting, more duplicates (idempotency is critical)
