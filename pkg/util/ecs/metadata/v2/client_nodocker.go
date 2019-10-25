// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

// +build !docker

package v2

type Client struct{}

func NewDefaultClient() *Client {
	return new(Client)
}

func (c *Client) GetTask() (*Task, error) {
	return new(Task), nil
}

func (c *Client) GetTaskWithTags() (*Task, error) {
	return new(Task), nil
}
