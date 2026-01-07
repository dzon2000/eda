# 0002 - Add `eventId` Field to `OrderCreated` Schema

- **Status:** Proposed <!-- or Accepted / Superseded -->
- **Date:** 2026-01-07

## Context

We emit an `OrderCreated` event using an Avro schema at:

- `schemas/order-created.avsc`
- Namespace: `io.pw.orders.v1`
- Record name: `OrderCreated`

Current relevant schema excerpt (after ADR-0001):

```json
{
    "type": "record",
    "name": "OrderCreated",
    "namespace": "io.pw.orders.v1",
    "fields": [
        { "name": "orderId", "type": "string" },
        { "name": "customerId", "type": "string" },
        { "name": "amount", "type": "double" },
        { "name": "createdAt", "type": "string" },
        { "name": "discount", "type": ["null", "double"], "default": null }
    ]
}
```

Currently, consumers cannot reliably implement idempotency because there is no stable, event-level identifier. Using `orderId` is not sufficient, as multiple events may exist per order in the future, and retries of the same event should be distinguishable from new events.

## Decision

We will add a required `eventId` field to the `OrderCreated` schema with a default value to preserve backward compatibility:

```json
{
    "name": "eventId",
    "type": "string",
    "default": ""
}
```

Resulting schema excerpt:

```json
{
    "type": "record",
    "name": "OrderCreated",
    "namespace": "io.pw.orders.v1",
    "fields": [
        { "name": "eventId", "type": "string", "default": "" },
        { "name": "orderId", "type": "string" },
        { "name": "customerId", "type": "string" },
        { "name": "amount", "type": "double" },
        { "name": "createdAt", "type": "string" },
        { "name": "discount", "type": ["null", "double"], "default": null }
    ]
}
```

Notes:

- `eventId` is a unique identifier for a specific event instance, not for the aggregate (`orderId`).
- Producers will generate a unique `eventId` (UUID) for every emitted `OrderCreated` event.
- The empty-string default keeps the schema backward compatible for events written before this ADR.

## Consequences

- **Positive:**
    - Enables consumers to implement idempotent processing (e.g., deduplication based on `eventId`).
    - Facilitates tracing and debugging across services by referencing a single stable event identifier.
    - Backward compatible: older events that lack `eventId` will deserialize with the default value.

- **Negative:**
    - Requires all producers of `OrderCreated` to start generating and populating `eventId`.
    - Consumers must treat empty `eventId` as “legacy/no-idempotency” and handle that case explicitly.

## Implementation Notes

- Update `schemas/order-created.avsc` to include the `eventId` field (preferably first in the `fields` array).
- Update producer code:
    - Extend `OrderCreated` struct with `EventID string`.
    - Populate `eventId` in `ToMap()`:

        ```go
        data := map[string]interface{}{
                "eventId":    o.EventID,
                "orderId":    o.OrderID,
                "customerId": o.CustomerID,
                "amount":     o.Amount,
                "createdAt":  o.CreatedAt,
                // discount as in ADR-0001
        }
        ```

    - Ensure `NewOrderCreatedEvent` (or equivalent factory) generates a unique `EventID` per event.
- Consumers:
    - Optionally introduce an idempotency store keyed by `eventId`.
    - Treat missing/empty `eventId` as non-idempotent legacy events.
- Communicate the new field and its semantics (uniqueness, immutability, usage for idempotency) to all teams consuming `OrderCreated`.
- Register the updated schema in Schema Registry under the existing subject (`orders.v1-value`) and verify `BACKWARD` compatibility. 