// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package metadata

import (
	"sync"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/DataDog/datadog-agent/pkg/util/retry"

	v1 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v1"
	v2 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v2"
	v3 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v3"
)

var globalUtil util

type util struct {
	// used to setup the ECSUtil
	initRetryV1 retry.Retrier
	initRetryV3 retry.Retrier
	initV1      sync.Once
	initV2      sync.Once
	initV3      sync.Once
	v1          *v1.Client
	v2          *v2.Client
	v3          *v3.Client
}

func V1() (*v1.Client, error) {
	globalUtil.initV1.Do(func() {
		globalUtil.initRetryV1.SetupRetrier(&retry.Config{
			Name:          "ecsutil-meta-v1",
			AttemptMethod: initV1,
			Strategy:      retry.RetryCount,
			RetryCount:    10,
			RetryDelay:    30 * time.Second,
		})
	})
	if err := globalUtil.initRetryV1.TriggerRetry(); err != nil {
		log.Debugf("ECS metadata v1 client init error: %s", err)
		return nil, err
	}
	return globalUtil.v1, nil
}

func V2() *v2.Client {
	globalUtil.initV2.Do(func() {
		globalUtil.v2 = v2.NewDefaultClient()
	})
	return globalUtil.v2
}

func V3(containerID string) (*v3.Client, error) {
	return v3.NewClientForContainer(containerID)
}

func V3FromCurrentTask() (*v3.Client, error) {
	globalUtil.initV3.Do(func() {
		globalUtil.initRetryV3.SetupRetrier(&retry.Config{
			Name:          "ecsutil-meta-v3",
			AttemptMethod: initV3,
			Strategy:      retry.RetryCount,
			RetryCount:    10,
			RetryDelay:    30 * time.Second,
		})
	})
	if err := globalUtil.initRetryV3.TriggerRetry(); err != nil {
		log.Debugf("ECS metadata v3 client init error: %s", err)
		return nil, err
	}
	return globalUtil.v3, nil
}

func initV1() error {
	client, err := v1.NewAutodetectedClient()
	if err != nil {
		return err
	}
	globalUtil.v1 = client
	return nil
}

func initV3() error {
	client, err := v3.NewClientForCurrentTask()
	if err != nil {
		return err
	}
	globalUtil.v3 = client
	return nil
}
