// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build docker

package collectors

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-agent/pkg/tagger/utils"
	"github.com/DataDog/datadog-agent/pkg/util/docker"
	ecsutil "github.com/DataDog/datadog-agent/pkg/util/ecs"
	ecsmeta "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata"
	v1 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v1"
	v3 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v3"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

func (c *ECSCollector) parseTasks(tasks []v1.Task, targetDockerID string) ([]*TagInfo, error) {
	var output []*TagInfo
	now := time.Now()
	for _, task := range tasks {
		// We only want to collect tasks without a STOPPED status.
		if task.KnownStatus == "STOPPED" {
			continue
		}
		for _, container := range task.Containers {
			// Only collect new containers + the targeted container, to avoid empty tags on race conditions
			if c.expire.Update(container.DockerID, now) || container.DockerID == targetDockerID {
				tags := utils.NewTagList()
				tags.AddLow("task_version", task.Version)
				tags.AddLow("task_name", task.Family)
				tags.AddLow("task_family", task.Family)
				tags.AddLow("ecs_container_name", container.Name)

				if c.clusterName != "" {
					tags.AddLow("cluster_name", c.clusterName)
				}

				if ecsutil.HasECSResourceTags() {
					if task, err := fetchTaskWithTagsV3(container.DockerID); err != nil {
						log.Warnf("Unable to get resource tags for container %s: %s", container.DockerID, err)
					} else {
						addResourceTags(tags, task.ContainerInstanceTags)
						addResourceTags(tags, task.TaskTags)
					}
				}

				tags.AddOrchestrator("task_arn", task.Arn)

				low, orch, high := tags.Compute()

				info := &TagInfo{
					Source:               ecsCollectorName,
					Entity:               docker.ContainerIDToTaggerEntityName(container.DockerID),
					HighCardTags:         high,
					OrchestratorCardTags: orch,
					LowCardTags:          low,
				}
				output = append(output, info)
			}
		}
	}
	return output, nil
}

func fetchTaskWithTagsV3(containerID string) (*v3.Task, error) {
	metaV3, err := ecsmeta.V3(containerID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize client for metadata v3 API: %s", err)
	}
	task, err := metaV3.GetTaskWithTags()
	if err != nil {
		return nil, fmt.Errorf("failed to get task with tags from metadata v3 API: %s", err)
	}
	return task, nil
}
