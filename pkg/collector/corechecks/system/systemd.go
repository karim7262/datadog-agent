// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package system

import (
	"fmt"

	"github.com/DataDog/datadog-agent/pkg/aggregator"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/coreos/go-systemd/dbus"

	core "github.com/DataDog/datadog-agent/pkg/collector/corechecks"
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

	conn, err := dbus.New()
	if err != nil {
		log.Error("New Connection: ", err)
		return nil
	}
	defer conn.Close()

	units, err := conn.ListUnits()

	activeUnitCounter := 0
	allUnitCounter := 0
	for _, unit := range units {
		log.Debugf("[unit] %s: ActiveState=%s, SubState=%s", unit.Name, unit.ActiveState, unit.SubState)
		if unit.ActiveState == "active" {
			activeUnitCounter++
		}
		allUnitCounter++
	}
	sender.Gauge("test.systemd.unit.active.count", float64(activeUnitCounter), "", nil)
	sender.Gauge("test.systemd.unit.all.count", float64(allUnitCounter), "", nil)

	for _, unit := range units {
		log.Debugf("[unit] %s: ActiveState=%s, SubState=%s", unit.Name, unit.ActiveState, unit.SubState)
		if unit.ActiveState == "active" {

			tags := []string{fmt.Sprintf("unit_name:%s", unit.Name)}

			cpuProperty, err := conn.GetServiceProperty(unit.Name, "CPUUsageNSec")
			if err != nil {
				log.Error("New Connection: ", err)
			} else {
				sender.Gauge("test.systemd.unit.cpu", float64(cpuProperty.Value.Value().(uint64)), "", tags)
			}
			memoryroperty, err := conn.GetServiceProperty(unit.Name, "MemoryCurrent")
			if err != nil {
				log.Error("New Connection: ", err)
			} else {
				sender.Gauge("test.systemd.unit.memory", float64(memoryroperty.Value.Value().(uint64)), "", tags)
			}
			tasksProperty, err := conn.GetServiceProperty(unit.Name, "TasksCurrent")
			if err != nil {
				log.Error("New Connection: ", err)
			} else {
				sender.Gauge("test.systemd.unit.tasks", float64(tasksProperty.Value.Value().(uint64)), "", tags)
			}
		}
	}

	fmt.Println("==============\n")

	p, err := conn.GetUnitProperties("sshd.service")

	if err != nil {
		fmt.Println("GetUnitProperties: ", err)
		return nil
	}

	for k, v := range p {
		fmt.Printf("%50v >>> %v\n", k, v)
	}

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
