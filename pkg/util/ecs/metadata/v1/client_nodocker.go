// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

// +build !docker

package v1

import "github.com/DataDog/datadog-agent/pkg/util/docker"

// Client represents a client for the metadata v1 API endpoint.
type Client struct{}

// NewAutodetectedClient detects the metadata v1 API endpoint and creates a new
// client for it. Returns an error if it was not possible to find the endpoint.
func NewAutodetectedClient() (*Client, error) {
	return nil, docker.ErrDockerNotCompiled
}
