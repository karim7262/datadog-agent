// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

// +build docker

package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/docker"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	defaultAgentPort = 51678
)

func detectAgentURL() (string, error) {
	urls := make([]string, 0, 3)

	if len(config.Datadog.GetString("ecs_agent_url")) > 0 {
		urls = append(urls, config.Datadog.GetString("ecs_agent_url"))
	}

	if config.IsContainerized() {
		// List all interfaces for the ecs-agent container
		agentURLS, err := getAgentContainerURLS()
		if err != nil {
			log.Debugf("could not inspect ecs-agent container: %s", err)
		} else {
			urls = append(urls, agentURLS...)
		}
		// Try the default gateway
		gw, err := docker.DefaultGateway()
		if err != nil {
			log.Debugf("could not get docker default gateway: %s", err)
		}
		if gw != nil {
			urls = append(urls, fmt.Sprintf("http://%s:%d/", gw.String(), defaultAgentPort))
		}

	}

	// Always try the localhost URL.
	urls = append(urls, fmt.Sprintf("http://localhost:%d/", defaultAgentPort))

	// Try the default IP for awsvpc mode
	urls = append(urls, fmt.Sprintf("http://169.254.172.1:%d/", defaultAgentPort))

	detected := testURLs(urls, 1*time.Second)
	if detected != "" {
		return detected, nil
	}

	return "", fmt.Errorf("could not detect ECS agent, tried URLs: %s", urls)
}

func getAgentContainerURLS() ([]string, error) {
	var urls []string

	du, err := docker.GetDockerUtil()
	if err != nil {
		return nil, err
	}
	ecsConfig, err := du.Inspect(config.Datadog.GetString("ecs_agent_container_name"), false)
	if err != nil {
		return nil, err
	}

	for _, network := range ecsConfig.NetworkSettings.Networks {
		ip := network.IPAddress
		if ip != "" {
			urls = append(urls, fmt.Sprintf("http://%s:%d/", ip, defaultAgentPort))
		}
	}

	// Add the container hostname, as it holds the instance's private IP when ecs-agent
	// runs in the (default) host network mode. This allows us to connect back to it
	// from an agent container running in awsvpc mode.
	if ecsConfig.Config != nil && ecsConfig.Config.Hostname != "" {
		urls = append(urls, fmt.Sprintf("http://%s:%d/", ecsConfig.Config.Hostname, defaultAgentPort))
	}

	return urls, nil
}

// testURLs trys a set of URLs and returns the first one that succeeds.
func testURLs(urls []string, timeout time.Duration) string {
	client := &http.Client{Timeout: timeout}
	for _, url := range urls {
		r, err := client.Get(url)
		if err != nil {
			continue
		}
		if r.StatusCode != http.StatusOK {
			continue
		}
		var resp Commands
		if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
			fmt.Printf("decode err: %s\n", err)
			continue
		}
		if len(resp.AvailableCommands) > 0 {
			return url
		}
	}
	return ""
}
