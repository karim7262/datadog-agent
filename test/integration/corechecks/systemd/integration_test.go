// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package systemd

//func init() {
//	registerComposeFile("basemetrics.compose")
//}
//
//func TestContainerMetricsTagging(t *testing.T) {
//
//	expectedMetrics := map[string][]string{
//		"Gauge": {
//			"systemd.unit.count",
//			"systemd.unit.loaded.count",
//
//			"systemd.unit.uptime",
//			"systemd.unit.loaded",
//			"systemd.unit.active",
//
//			"systemd.socket.n_connections",
//			"systemd.socket.n_accepted",
//
//			//"systemd.service.memory_current",
//			//"systemd.service.tasks_current",
//
//			// centos/systemd:latest contains systemd v219, it does not contain CPUUsageNSec and NRestarts yet
//			// "systemd.service.cpu_usage_n_sec",
//			// "systemd.service.n_restarts",
//
//			"systemd.unit.count",
//			"systemd.unit.loaded.count",
//		},
//	}
//
//	for method, metricList := range expectedMetrics {
//		for _, metric := range metricList {
//			sender.AssertCalled(t, method, metric, mock.Anything, "", mock.Anything)
//		}
//	}
//
//}
