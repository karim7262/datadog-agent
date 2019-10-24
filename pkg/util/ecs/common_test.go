// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build docker

package ecs

import (
	"net"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/util/containers"
	v2 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v2"
	"github.com/stretchr/testify/assert"
)

func TestParseContainerNetworkAddresses(t *testing.T) {
	ports := []v2.Port{
		{
			ContainerPort: 80,
			Protocol:      "tcp",
		},
		{
			ContainerPort: 7000,
			Protocol:      "udp",
		},
	}
	networks := []v2.Network{
		{
			NetworkMode:   "awsvpc",
			IPv4Addresses: []string{"10.0.2.106"},
		},
		{
			NetworkMode:   "awsvpc",
			IPv4Addresses: []string{"10.0.2.107"},
		},
	}
	expectedOutput := []containers.NetworkAddress{
		{
			IP:       net.ParseIP("10.0.2.106"),
			Port:     80,
			Protocol: "tcp",
		},
		{
			IP:       net.ParseIP("10.0.2.106"),
			Port:     7000,
			Protocol: "udp",
		},
		{
			IP:       net.ParseIP("10.0.2.107"),
			Port:     80,
			Protocol: "tcp",
		},
		{
			IP:       net.ParseIP("10.0.2.107"),
			Port:     7000,
			Protocol: "udp",
		},
	}
	result := parseContainerNetworkAddresses(ports, networks, "mycontainer")
	assert.Equal(t, expectedOutput, result)
}
