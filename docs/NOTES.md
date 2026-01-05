# What we are building

This is EDA (event driven architecture) project to learn how events are transmitted in distributed systems.

## Components

The project will use Kafka for event distribution and have Golang based producers and consumers. For schema versioning we will use 
Schema Registry by Confluent.

## Schema Registry

- Supports versioning of event schemas

## AVRO

- Serialization/Deserialization of complex data.
- Schema-based serialization
- Can use a newer backward-compatible schema to decode (deserialize) data that was encoded (serialized) with an older schema