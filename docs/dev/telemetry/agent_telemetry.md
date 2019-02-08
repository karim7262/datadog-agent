# Agent telemetry

## Metric pipeline

Here is a high level diagram of the agent metric pipeline:

![agent-diagram](./agent.svg)

## DogStatsD server

The DogStatsD (aka DSD) server task is to process incoming [dogstatsd] packet/datagrams. They are then forwarded to the aggregator.

We collect the following metrics:

### Dogstatsd Intake

- `datadog.agent.dogstatsd.packet_in`
- `datadog.agent.dogstatsd.packet_in_errors`
