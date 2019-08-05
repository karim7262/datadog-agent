// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package docker

import (
	"testing"
)

func init() {
	registerComposeFile("basemetrics.compose")
}

func TestContainerMetricsTagging(t *testing.T) {

	//expectedTags := []string{
	//	"test:e2e",
	//}

	//sender.AssertCalled(t, "Gauge", "systemd.unit.loaded.count", 267, "", [])
	sender.AssertCalled(t, "Gauge", "systemd.unit.loaded.count", float64(267), "", []string(nil))

	//expectedMetrics := map[string][]string{
	//	"Gauge": {
	//		"docker.cpu.shares",
	//		"docker.kmem.usage",
	//		"docker.mem.cache",
	//		"docker.mem.rss",
	//		"docker.mem.in_use",
	//		"docker.mem.limit",
	//		"docker.mem.failed_count",
	//		"docker.mem.soft_limit",
	//		"docker.container.size_rwx",
	//		"docker.container.size_rootfs",
	//		"docker.thread.count",
	//	},
	//	"Rate": {
	//		"docker.cpu.system",
	//		"docker.cpu.user",
	//		"docker.cpu.usage",
	//		"docker.cpu.throttled",
	//		"docker.io.read_bytes",
	//		"docker.io.write_bytes",
	//		"docker.net.bytes_sent",
	//		"docker.net.bytes_rcvd",
	//	},
	//}

	//tags := []string{
	//	"test:e2e",
	//}
	//
	//ok := sender.AssertMetricTaggedWith(t, "Gauge", "docker.containers.running", tags)
	//if !ok {
	//	log.Warnf("Missing Gauge docker.containers.running with tags %s", tags)
	//}
	//
	//
	//for method, metricList := range expectedMetrics {
	//	for _, metric := range metricList {
	//		ok := sender.AssertMetricTaggedWith(t, method, metric, expectedTags)
	//		if !ok {
	//			log.Warnf("Missing %s %s with tags %s", method, metric, expectedTags)
	//		}
	//
	//		// Excluded pause container
	//		ok = sender.AssertMetricNotTaggedWith(t, method, metric, pauseTags)
	//		if !ok {
	//			log.Warnf("Shouldn't call %s %s with tags %s", method, metric, pauseTags)
	//		}
	//	}
	//}
}
