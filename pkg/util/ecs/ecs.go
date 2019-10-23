// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build docker

package ecs

import (
	"os"
	"sync"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/util/cache"
	"github.com/DataDog/datadog-agent/pkg/util/retry"

	v1 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v1"
	v2 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v2"
	v3 "github.com/DataDog/datadog-agent/pkg/util/ecs/metadata/v3"
)

const (
	// Cache the fact we're running on ECS Fargate
	isFargateInstanceCacheKey = "IsFargateInstanceCacheKey"
	// Cache the fact resources tags are exposed on ECS Fargate
	hasFargateResourceTagsCacheKey = "HasFargateResourceTagsCacheKey"
	// Cache the fact resources tags are exposed on ECS over EC2
	hasECSResourceTagsCacheKey = "HasECSResourceTagsCacheKey"
	// CloudProviderName contains the inventory name of for ECS
	CloudProviderName = "AWS"
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

// IsRunningOn returns true if the agent is running on ECS/Fargate
func IsRunningOn() bool {
	return IsECSInstance() || IsFargateInstance()
}

func MetaV1() (*v1.Client, error) {
	globalUtil.initV1.Do(func() {
		globalUtil.initRetryV1.SetupRetrier(&retry.Config{
			Name:          "ecsmetav1",
			AttemptMethod: initV1,
			Strategy:      retry.RetryCount,
			RetryCount:    10,
			RetryDelay:    30 * time.Second,
		})
	})
	if err := globalUtil.initRetryV1.TriggerRetry(); err != nil {
		log.Debugf("ECS init error: %s", err)
		return nil, err
	}
	return globalUtil.v1, nil
}

func MetaV2() *v2.Client {
	globalUtil.initV2.Do(func() {
		globalUtil.v2 = v2.NewDefaultClient()
	})
	return globalUtil.v2
}

func MetaV3(containerID string) (*v3.Client, error) {
	return v3.NewClientForContainer(containerID)
}

func MetaV3InCurrentTask() (*v3.Client, error) {
	globalUtil.initV3.Do(func() {
		globalUtil.initRetryV3.SetupRetrier(&retry.Config{
			Name:          "ecsmetav3",
			AttemptMethod: initV3,
			Strategy:      retry.RetryCount,
			RetryCount:    10,
			RetryDelay:    30 * time.Second,
		})
	})
	if err := globalUtil.initRetryV3.TriggerRetry(); err != nil {
		log.Debugf("ECS init error: %s", err)
		return nil, err
	}
	return globalUtil.v3, nil
}

// IsECSInstance returns whether the agent is running in ECS.
func IsECSInstance() bool {
	_, err := MetaV1()
	return err == nil
}

// IsFargateInstance returns whether the agent is in an ECS fargate task.
// It detects it by getting and unmarshalling the metadata API response.
func IsFargateInstance() bool {
	return cacheQueryBool(isFargateInstanceCacheKey, func() (bool, time.Duration) {

		// This envvar is set to AWS_ECS_EC2 on classic EC2 instances
		// Versions 1.0.0 to 1.3.0 (latest at the time) of the Fargate
		// platform set this envvar.
		// If Fargate detection were to fail, running a container with
		// `env` as cmd will allow to check if it is still present.
		if os.Getenv("AWS_EXECUTION_ENV") != "AWS_ECS_FARGATE" {
			return newBoolEntry(false)
		}

		_, err := MetaV2().GetTask()
		if err != nil {
			log.Debug(err)
			return newBoolEntry(false)
		}

		return newBoolEntry(true)
	})
}

// HasFargateResourceTags returns whether the metadata endpoint in Fargate
// exposes resource tags.
func HasFargateResourceTags() bool {
	return cacheQueryBool(hasFargateResourceTagsCacheKey, func() (bool, time.Duration) {
		_, err := MetaV2().GetTaskWithTags()
		return newBoolEntry(err == nil)
	})
}

// HasECSResourceTags returns whether the metadata endpoint in ECS exposes
// resource tags.
func HasECSResourceTags() bool {
	return cacheQueryBool(hasECSResourceTagsCacheKey, func() (bool, time.Duration) {
		client, err := MetaV3InCurrentTask()
		if err != nil {
			return newBoolEntry(false)
		}
		_, err = client.GetTaskWithTags()
		return newBoolEntry(err == nil)
	})
}

func cacheQueryBool(cacheKey string, cacheMissEvalFunc func() (bool, time.Duration)) bool {
	if cachedValue, found := cache.Cache.Get(cacheKey); found {
		if v, ok := cachedValue.(bool); ok {
			return v
		}
		log.Errorf("Invalid cache format for key %q: forcing a cache miss", cacheKey)
	}

	newValue, ttl := cacheMissEvalFunc()
	cache.Cache.Set(cacheKey, newValue, ttl)

	return newValue
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

func newBoolEntry(v bool) (bool, time.Duration) {
	if v == true {
		return v, 5 * time.Minute
	}
	return v, cache.NoExpiration
}
