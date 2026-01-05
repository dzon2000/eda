# 0001 - Add `discount` Field to `OrderCreated` Schema

- **Status:** Proposed <!-- or Accepted / Superseded -->
- **Date:** 2026-01-05

## Context

We emit an `OrderCreated` event using an Avro schema at:

- `schemas/order-created.avsc`
- Namespace: `io.pw.orders.v1`
- Record name: `OrderCreated`

Current relevant schema excerpt:

```json
{
    "type": "record",
    "name": "OrderCreated",
    "namespace": "io.pw.orders.v1",
    "fields": [
        { "name": "orderId", "type": "string" },
        { "name": "customerId", "type": "string" },
        { "name": "amount", "type": "double" },
        { "name": "createdAt", "type": "string" }
    ]
}
```

We need to represent discounts applied at order creation while preserving backward compatibility for existing consumers.

## Decision

We will add an optional `discount` field to the `OrderCreated` schema:

```json
{
    "name": "discount",
    "type": ["null", "double"],
    "default": null
}
```

Resulting schema excerpt:

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

Notes:

- The field is nullable with a default of `null` to maintain backward compatibility.
- Producers will populate `discount` for new use cases where a discount is applied.

## Consequences

- **Positive:**
    - Enables analytics and business logic that depend on discount information.
    - Backward compatible: existing consumers that ignore unknown fields continue to work.
- **Negative:**
    - Consumers that perform strict schema validation or use code generation must be updated to handle the new optional field.
    - Additional responsibility on producers to correctly calculate and populate `discount`.

## Implementation Notes

- Update `schemas/order-created.avsc` with the new field.
- Regenerate any schema-derived code (e.g., Avro classes) in producer and consumer services.
- Communicate the change and rollout plan to all teams consuming `OrderCreated`.