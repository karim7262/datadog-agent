// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

// +build docker

package v3

import (
	"fmt"
	"os"
	"strings"

	"github.com/DataDog/datadog-agent/pkg/util/docker"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	defaultMetadataURIEnvVariable = "ECS_CONTAINER_METADATA_URI"
)

func getAgentURLFromEnv() (string, error) {
	agentURL, found := os.LookupEnv(defaultMetadataURIEnvVariable)
	if !found {
		return "", fmt.Errorf("Could not initialize client: missing metadata v3 URL")
	}
	return agentURL, nil
}

func getAgentURLFromDocker(containerID string) (string, error) {
	du, err := docker.GetDockerUtil()
	if err != nil {
		return "", err
	}

	container, err := du.Inspect(containerID, false)
	if err != nil {
		return "", err
	}

	for _, env := range container.Config.Env {
		substrings := strings.Split(env, "=")
		if len(substrings) != 2 {
			log.Tracef("invalid container env format: %s", env)
		}

		k := substrings[0]
		v := substrings[1]

		if k == defaultMetadataURIEnvVariable {
			return v, nil
		}
	}

	return "", fmt.Errorf("metadata v3 URL not found in container %s", containerID)
}
