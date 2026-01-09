# Consumer correctness & rebalancing

## Goal

- Understand why consumers stall or duplicate work
- Control rebalancing behavior
- Design consumers that survive restarts and scaling

What we will learn:
- Why only one consumer reads a partition
- Why idle consumers are normal
- Why scaling consumers ≠ scaling throughput

## Consumer groups

A consumer group is:
- A named group of consumers
- Cooperating to read one topic
- Each partition is assigned to exactly one consumer in the group

> Kafka scales consumers by partitions, not by consumer count.

Kafka remembers:
- Last committed offset
- Per partition
- Per group

Change GroupID → Kafka treats you as a new application

This is why:
- Messages “reappear”
- Or “disappear”

Offsets belong to group.id, topic, partition

Not:
- Consumer instance
- Host
- Pod
- Container

This enables:
- Restart
- Scaling
- Failover

What happens on startup
- Consumer joins group
- Group coordinator assigns partitions
- Consumer starts from committed offset
- Poll loop begins

**Important**

1. Consumer count ≤ partition count
2. GroupID change = new application
3. Rebalances are normal
4. Slow consumers cause rebalances
5. Offsets are group-scoped state