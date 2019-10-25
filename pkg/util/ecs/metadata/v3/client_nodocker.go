// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc

// +build !docker

package v3

import "github.com/DataDog/datadog-agent/pkg/util/docker"

// Client represents a client for a metadata v3 API endpoint.
type Client struct{}

// NewClientForCurrentTask detects the metadata API v3 endpoint from the current
// task and creates a new client for it.
func NewClientForCurrentTask() (*Client, error) {
	return nil, docker.ErrDockerNotCompiled
}

// NewClientForContainer detects the metadata API v3 endpoint for the specified
// container and creates a new client for it.
func NewClientForContainer(containerID string) (*Client, error) {
	return nil, docker.ErrDockerNotCompiled
}
