// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

// +build docker

package v3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const (
	// Metadata v3 API paths
	taskMetadataPath         = "/task"
	taskMetadataWithTagsPath = "/taskWithTags"
	containerStatsPath       = "/stats/"

	// Default client configuration
	endpointTimeout = 500 * time.Millisecond
)

// Client represents a client for a metadata v3 API endpoint.
type Client struct {
	agentURL string
}

// NewClient creates a new client for the specified metadata v3 API endpoint.
func NewClient(agentURL string) *Client {
	if strings.HasSuffix(agentURL, "/") {
		agentURL = strings.TrimSuffix(agentURL, "/")
	}
	return &Client{
		agentURL: agentURL,
	}
}

// NewClientForCurrentTask detects the metadata API v3 endpoint from the current
// task and creates a new client for it.
func NewClientForCurrentTask() (*Client, error) {
	agentURL, err := getAgentURLFromEnv()
	if err != nil {
		return nil, err
	}
	return NewClient(agentURL), nil
}

// NewClientForContainer detects the metadata API v3 endpoint for the specified
// container and creates a new client for it.
func NewClientForContainer(id string) (*Client, error) {
	agentURL, err := getAgentURLFromDocker(id)
	if err != nil {
		return nil, err
	}
	return NewClient(agentURL), nil
}

// GetTask returns the current task.
func (c *Client) GetTask() (*Task, error) {
	return c.getTaskMetadataAtPath(taskMetadataPath)
}

// GetTaskWithTags returns the current task, including propagated resource tags.
func (c *Client) GetTaskWithTags() (*Task, error) {
	return c.getTaskMetadataAtPath(taskMetadataWithTagsPath)
}

func (c *Client) get(path string, v interface{}) error {
	client := http.Client{Timeout: endpointTimeout}
	url := c.makeURL(path)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected HTTP status code in metadata v3 reply: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("Failed to decode metadata v3 JSON payload to type %s: %s", reflect.TypeOf(v), err)
	}

	return nil
}

func (c *Client) getTaskMetadataAtPath(path string) (*Task, error) {
	var t Task
	if err := c.get(path, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) makeURL(path string) string {
	return fmt.Sprintf("%s%s", c.agentURL, path)
}
