// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

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
	// ECS agent defaults
	defaultAgentURL = "http://169.254.170.2/"

	// Metadata v2 API paths
	taskMetadataPath         = "/task"
	taskMetadataWithTagsPath = "/taskWithTags"
	containerStatsPath       = "/stats/"

	// Default client configuration
	endpointTimeout = 500 * time.Millisecond
)

type Client struct {
	agentURL string
}

func NewClient(agentURL string) *Client {
	if !strings.HasSuffix(agentURL, "/") {
		agentURL += "/"
	}
	return &Client{
		agentURL: agentURL,
	}
}

func NewClientForCurrentTask() (*Client, error) {
	agentURL, err := getAgentURLFromEnv()
	if err != nil {
		return nil, err
	}
	return NewClient(agentURL), nil
}

func NewClientForContainer(id string) (*Client, error) {
	agentURL, err := getAgentURLFromDocker(id)
	if err != nil {
		return nil, err
	}
	return NewClient(agentURL), nil
}

func (c *Client) GetTask() (*Task, error) {
	return c.getTaskMetadataAtPath(taskMetadataPath)
}

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
		return fmt.Errorf("Failed decoding metadata v3 JSON object to type %s - %s", reflect.TypeOf(v), err)
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
	return fmt.Sprintf("%sv3%s", c.agentURL, path)
}
