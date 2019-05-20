// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.
// +build !windows

package systemd

import (
	"fmt"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	"github.com/coreos/go-systemd/dbus"
	"github.com/shirou/gopsutil/cpu"
)

var (
	firstSample = []cpu.TimesStat{
		{
			CPU:       "cpu-total",
			User:      1229386,
			Nice:      623,
			System:    263584,
			Idle:      25496761,
			Iowait:    12113,
			Irq:       10,
			Softirq:   1151,
			Steal:     0,
			Guest:     0,
			GuestNice: 0,
			Stolen:    0,
		},
	}
	secondSample = []cpu.TimesStat{
		{
			CPU:       "cpu-total",
			User:      1229586,
			Nice:      625,
			System:    268584,
			Idle:      25596761,
			Iowait:    12153,
			Irq:       15,
			Softirq:   1451,
			Steal:     2,
			Guest:     0,
			GuestNice: 0,
			Stolen:    0,
		},
	}
)

var sample = firstSample

func CPUTimes(bool) ([]cpu.TimesStat, error) {
	return sample, nil
}

// // Conn is a connection to systemd's dbus endpoint.
// type Conn struct {
// }

// func (c *Conn) ListUnits() ([]coredbus.UnitStatus, error) {
// 	return nil, nil
// }

// Conn is a connection to systemd's dbus endpoint.
// type Conn struct {
// }

// type MyMockedObject struct {
// 	Conn
// }

type Conn struct {
}

func dbusNewFake() (*dbus.Conn, error) {
	fmt.Println("my dbusNewFake")
	return nil, nil
}

func connListUnitsFake(c *dbus.Conn) ([]dbus.UnitStatus, error) {
	fmt.Println("my connListUnitsFake")
	return []dbus.UnitStatus{
		{Name: "unit1", ActiveState: "active"},
		{Name: "unit2", ActiveState: "active"},
		{Name: "unit3", ActiveState: "inactive"},
	}, nil
}

func connCloseFake(c *dbus.Conn) {
}

func TestSystemdCheckLinux(t *testing.T) {
	dbusNew = dbusNewFake
	connListUnits = connListUnitsFake
	connClose = connCloseFake

	// create an instance of our test object
	systemdCheck := new(SystemdCheck)
	systemdCheck.Configure(nil, nil)

	// setup expectations
	mock := mocksender.NewMockSender(systemdCheck.ID())

	mock.On("Gauge", "systemd.unit.cpu", 1.0, "", []string(nil)).Return().Times(1)
	mock.On("Gauge", "systemd.unit.active.count", 2.0, "", []string(nil)).Return().Times(1)
	mock.On("Commit").Return().Times(1)

	systemdCheck.Run()

	mock.AssertExpectations(t)
	mock.AssertNumberOfCalls(t, "Gauge", 2)
	mock.AssertNumberOfCalls(t, "Commit", 1)

}
