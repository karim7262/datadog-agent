# Agent telemetry

## Metric pipeline

Here is a high level diagram of the agent metric pipeline:

![agent-diagram](./agent.svg)

## DogStatsD server

The DogStatsD (aka DSD) server task is to process incoming dogstatsd datagrams sent by dogstatsd clients through UDP and UDS. They are then forwarded to the aggregator.

We collect the following metrics:

#### `datadog.agent.dogstatsd.udp.datagrams_in` - *count*

Total number of datagrams received by the UDP listener.

#### `datadog.agent.dogstatsd.udp.read_errors` - *count*

Total number of errors encountered while calling `UDPConn.ReadFrom`.
Note: this is not the total number of dropped datagrams.

#### `datadog.agent.dogstatsd.uds.datagrams_in` - *count*

Total number of datagrams received by the UDS listener.

#### `datadog.agent.dogstatsd.uds.read_errors` - *count*

Total number of errors encountered while calling `UDSConn.ReadFrom`.
Note: this is not the total number of dropped datagrams.

#### `datadog.agent.dogstatsd.uds.origin_detection_errors` - *count*

Total number of errors encountered while performing origin detection.

#### `datadog.agent.dogstatsd.parser.messages` - *count*

Total number of "messages" extracted from raw datagrams. A message is one line of text from the datagram and could be a metric, a service check, an event or something invalid.

#### `datadog.agent.dogstatsd.parser.errors` - *count*

Total number of errors encountered while parsing "messages".
Note: This total canbe supperior to `metric_errors` + `event_errors` + `service_check_errors` because it also includes messages that the parser could not identify.

#### `datadog.agent.dogstatsd.parser.metrics` - *count*

Total number of messages that were identified as metrics. This includes the ones that were not sucessfully parsed.

#### `datadog.agent.dogstatsd.parser.metrics_errors` - *count*

Total number of errors encoutered while parsing a metric message.

#### `datadog.agent.dogstatsd.parser.service_checks` - *count*

Total number of messages that were identified as service checks. This includes the ones that were not sucessfully parsed.

#### `datadog.agent.dogstatsd.parser.service_checks_errors` - *count*

Total number of errors encoutered while parsing a service check message.

#### `datadog.agent.dogstatsd.parser.events` - *count*

Total number of messages that were identified as events. This includes the ones that were not sucessfully parsed.

#### `datadog.agent.dogstatsd.parser.events_errors` - *count*

Total number of errors encoutered while parsing a event message.


## Checks Intake - AKA Sender

## Aggregator

## Serializer

## Forwarder
