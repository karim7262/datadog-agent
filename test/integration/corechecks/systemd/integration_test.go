// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package systemd

import (
	"github.com/DataDog/datadog-agent/pkg/aggregator/mocksender"
	"github.com/DataDog/datadog-agent/pkg/collector/check"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/systemd"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/DataDog/datadog-agent/test/integration/utils"
)

var sender *mocksender.MockSender
var systemdCheck check.Check

type SystemdSuite struct {
	suite.Suite
	compose       *utils.ComposeConf
	containerName string
}

// use a constructor to make the suite parametric
func NewSystemdSuite(containerName string) *SystemdSuite {
	return &SystemdSuite{
		containerName: containerName,
	}
}

func (suite *SystemdSuite) SetupSuite() {
	suite.compose = &utils.ComposeConf{
		ProjectName:         "systemd",
		FilePath:            "testdata/docker-compose.yaml",
		RemoveRebuildImages: true,
	}

	output, err := suite.compose.Start()
	require.NoError(suite.T(), err, string(output))
}

func (suite *SystemdSuite) TearDownSuite() {
	suite.compose.Stop()
}

// put configuration back in a known state before each test
func (suite *SystemdSuite) SetupTest() {
	// Setup check
	rawInstanceConfig := []byte(`
unit_names:
 - dbus.service
 - dbus.socket
system_bus_socket: /tmp/var/run/dbus/system_bus_socket
`)
	systemdCheck = systemd.SystemdFactory()
	systemdCheck.Configure(rawInstanceConfig, []byte(``))

	// Setup mock sender
	sender = mocksender.NewMockSender(systemdCheck.ID())
	sender.SetupAcceptAll()

	systemdCheck.Run()
}

func (suite *SystemdSuite) TestSystemd() {

	expectedMetrics := map[string][]string{
		"Gauge": {
			"systemd.unit.count",
			"systemd.unit.loaded.count",

			"systemd.unit.uptime",
			"systemd.unit.loaded",
			"systemd.unit.active",

			"systemd.socket.connection_count",
			"systemd.socket.connection_accepted_count",

			"systemd.service.memory_usage",
			"systemd.service.task_count",

			// centos/systemd:latest contains systemd v219, it does not contain CPUUsageNSec and NRestarts yet
			// "systemd.service.cpu_usage_n_sec",
			// "systemd.service.n_restarts",

			"systemd.unit.count",
			"systemd.unit.loaded.count",
		},
	}

	for method, metricList := range expectedMetrics {
		for _, metric := range metricList {
			sender.AssertCalled(suite.T(), method, metric, mock.Anything, "", mock.Anything)
		}
	}
}

func TestSystemdSuite(t *testing.T) {
	suite.Run(t, NewSystemdSuite("datadog-agent-test-systemd"))
}
