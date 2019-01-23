// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package telemetry

import (
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/dogstatsd/listeners"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

func collectDogstatsdTelemetry(sender aggregator.Sender) error {
	if !config.Datadog.GetBool("use_dogstatsd") {
		udpDrops, err := dogstatsdDropCount()
		if err != nil {
			log.Debugf("Could not get dogstatsd drop count: %s", err)
		} else {
			sender.MonotonicCount("datadog.agent.dogstatsd.udp.drops", float64(udpDrops), "", []string{})
		}
		sender.MonotonicCount("datadog.agent.telemetry.dogstatsd.udp.datagrams_total", float64(listeners.DogstatsdUDPDatagrams.Value()), "", []string{})
		sender.MonotonicCount("datadog.agent.telemetry.dogstatsd.udp.read_errors", float64(listeners.DogstatsdUDPDatagrams.Value()), "", []string{})
	}
	return nil
}
