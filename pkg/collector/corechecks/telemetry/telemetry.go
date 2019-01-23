// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package telemetry

import (
	"time"

	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/autodiscovery/integration"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
)

const telemetryCheckName = "telemetry"

// TelemetryCheck is a check collecting telemetry data about the agent
type telemetryCheck struct {
	core.CheckBase
	cfg            *telemetryConfig
	lastCollection time.Time
}

type telemetryInstanceConfig struct {
}

type telemetryInitConfig struct{}

type telemetryConfig struct {
	instance telemetryInstanceConfig
	initConf telemetryInitConfig
}

func (c *telemetryCheck) String() string {
	return "telemetry"
}

// Configure configure the data from the yaml
func (c *telemetryCheck) Configure(data integration.Data, initConfig integration.Data) error {
	err := c.CommonConfigure(data)
	if err != nil {
		return err
	}
	cfg := new(telemetryConfig)

	c.BuildID(data, initConfig)
	c.cfg = cfg

	return nil
}

// Run runs the check
func (c *telemetryCheck) Run() error {
	sender, err := aggregator.GetSender(c.ID())
	if err != nil {
		return err
	}

	err = collectDogstatsdTelemetry(sender)

	sender.Commit()

	return nil
}

func telemetryFactory() check.Check {
	return &telemetryCheck{
		CheckBase: core.NewCheckBase(telemetryCheckName),
	}
}

func init() {
	core.RegisterCheck(telemetryCheckName, telemetryFactory)
}
