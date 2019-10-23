// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build docker

package v1

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/cache"
	"github.com/DataDog/datadog-agent/pkg/util/docker"
)

func TestGetAgentContainerURLS(t *testing.T) {
	config.Datadog.SetDefault("ecs_agent_container_name", "ecs-agent-custom")
	defer config.Datadog.SetDefault("ecs_agent_container_name", "ecs-agent")

	// Setting mocked data in cache
	nets := make(map[string]*network.EndpointSettings)
	nets["bridge"] = &network.EndpointSettings{IPAddress: "172.17.0.2"}
	nets["foo"] = &network.EndpointSettings{IPAddress: "172.17.0.3"}

	co := types.ContainerJSON{
		Config: &container.Config{
			Hostname: "ip-172-29-167-5",
		},
		ContainerJSONBase: &types.ContainerJSONBase{},
		NetworkSettings: &types.NetworkSettings{
			Networks: nets,
		},
	}
	docker.EnableTestingMode()
	cacheKey := docker.GetInspectCacheKey("ecs-agent-custom", false)
	cache.Cache.Set(cacheKey, co, 10*time.Second)

	agentURLS, err := getAgentContainerURLS()
	assert.NoError(t, err)
	require.Len(t, agentURLS, 3)
	assert.Contains(t, agentURLS, "http://172.17.0.2:51678/")
	assert.Contains(t, agentURLS, "http://172.17.0.3:51678/")
	assert.Equal(t, "http://ip-172-29-167-5:51678/", agentURLS[2])
}
