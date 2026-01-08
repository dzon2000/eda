# Consumer correctness & rebalancing

## Goal

- Understand why consumers stall or duplicate work
- Control rebalancing behavior
- Design consumers that survive restarts and scaling

What we will learn:
- Why only one consumer reads a partition
- Why idle consumers are normal
- Why scaling consumers ≠ scaling throughput