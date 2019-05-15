// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package system

import (
	"fmt"
	"time"

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
		log.Error("New Connection Err: ", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnits()
	if err != nil {
		fmt.Println("ListUnits Err: ", err)
		return err
	}

	activeUnitCounter := 0
	for _, unit := range units {
		log.Debugf("[unit] %s: ActiveState=%s, SubState=%s", unit.Name, unit.ActiveState, unit.SubState)
		if unit.ActiveState == "active" {
			activeUnitCounter++
		}
	}

	sender.Gauge("test.systemd.unit.active.count", float64(activeUnitCounter), "", nil)

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
		}
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

// ==============================================
// May be used for for Service Checks and Events
// ==============================================
func sandboxEvent() {
	target := "graphical.target"

	log.Info("==> Sandbox sandboxEvent")

	conn, err := dbus.New()

	if err != nil {
		log.Error("New Err: ", err)
		return
	}

	err = conn.Subscribe()
	if err != nil {
		log.Error("Subscribe Err: ", err)
		return
	}

	evChan, errChan := conn.SubscribeUnits(time.Second)

	for {
		select {
		case changes := <-evChan:
			tCh, ok := changes[target]

			// Just continue until we see our event.
			if !ok {
				continue
			}

			log.Info("==> New ActiveState:", tCh.ActiveState)

			if tCh.ActiveState == "active" && tCh.Name == target {
				log.Error("Reached timeout")
			}
		case err = <-errChan:
			log.Error("change err: ", err)
		case <-time.After(10 * time.Second):
			log.Error("Reached timeout")
		}
	}

}
