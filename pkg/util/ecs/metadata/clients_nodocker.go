// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.

// +build !docker

package metadata

import (
	v1 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v1"
	v2 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v2"
	v3 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v3"
)

func V1() (*v1.Client, error) {
	return v1.NewAutodetectedClient()
}

func V2() *v2.Client {
	return v2.NewDefaultClient()
}

func V3(containerID string) (*v3.Client, error) {
	return v3.NewClientForContainer(containerID)
}

func V3FromCurrentTask() (*v3.Client, error) {
	return v3.NewClientForCurrentTask()
}
