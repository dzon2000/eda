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

