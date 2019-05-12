// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package system

import (
	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const systemdCheckName = "systemd"

// SystemdCheck doesn't need additional fields
type SystemdCheck struct {
	core.CheckBase
}

// Run executes the check
func (c *SystemdCheck) Run() error {
	sender, err := aggregator.GetSender(c.ID())
	if err != nil {
		return err
	}

	t, err := uptime()
	if err != nil {
		log.Errorf("system.SystemdCheck: could not retrieve uptime: %s", err)
		return err
	}
	log.Info("system.SystemdCheck: TEST")

	sender.Gauge("system.uptime", float64(t), "", nil)
	sender.Commit()

	return nil
}

func systemdFactory() check.Check {
	return &SystemdCheck{
		CheckBase: core.NewCheckBase(systemdCheckName),
	}
}

func init() {
	core.RegisterCheck(systemdCheckName, systemdFactory)
}
