// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// +build clusterchecks

package cloudfoundry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"code.cloudfoundry.org/locket/models"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/DataDog/datadog-agent/pkg/clusteragent/clusterchecks/types"
)

type testLocketClient struct {
	LockErr            error
	LockCalledTimes    int
	ReleaseErr         error
	ReleaseCalledTimes int
}

func (c *testLocketClient) Lock(ctx context.Context, in *models.LockRequest, opts ...grpc.CallOption) (*models.LockResponse, error) {
	c.LockCalledTimes++
	return &models.LockResponse{}, c.LockErr
}

func (c *testLocketClient) Fetch(ctx context.Context, in *models.FetchRequest, opts ...grpc.CallOption) (*models.FetchResponse, error) {
	panic("not implemented")
}

func (c *testLocketClient) Release(ctx context.Context, in *models.ReleaseRequest, opts ...grpc.CallOption) (*models.ReleaseResponse, error) {
	c.ReleaseCalledTimes++
	return &models.ReleaseResponse{}, c.ReleaseErr
}

func (c *testLocketClient) FetchAll(ctx context.Context, in *models.FetchAllRequest, opts ...grpc.CallOption) (*models.FetchAllResponse, error) {
	panic("not implemented")
}

func resetGlobalLeaderElector() {
	globalLeaderElector = nil
}

func TestLeaderElector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tc := &testLocketClient{}
	defer cancel()
	defer resetGlobalLeaderElector()
	leaderIP, err := RunGlobalLeaderElector(ctx, "url", "cafile", "certfile", "keyfile", "hostname", time.Millisecond, tc)
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 10)

	// when there's nil error returned by tc.Lock, we'll become leader
	leaderIPResult, err := leaderIP()
	assert.Nil(t, err)
	assert.Equal(t, types.LeaderIPResult{
		IP:           "",
		Redirectable: false,
	}, leaderIPResult)
	globalLeaderElector.Lock()
	assert.NotEqual(t, 0, tc.LockCalledTimes)
	assert.Equal(t, 0, tc.ReleaseCalledTimes)

	// when there's a non-nil error returned by tc.Lock, we'll become follower
	tc.LockCalledTimes = 0
	tc.LockErr = fmt.Errorf("Already locked, sorry")
	globalLeaderElector.Unlock()
	time.Sleep(time.Millisecond * 10)

	leaderIPResult, err = leaderIP()
	assert.Nil(t, err)
	assert.Equal(t, types.LeaderIPResult{
		IP:           NotALeaderResponse,
		Redirectable: false,
	}, leaderIPResult)
	globalLeaderElector.Lock()
	assert.NotEqual(t, 0, tc.LockCalledTimes)
	assert.Equal(t, 0, tc.ReleaseCalledTimes)

	// make sure we become leader again
	tc.LockErr = nil
	globalLeaderElector.Unlock()
	time.Sleep(time.Millisecond * 10)

	// mark context as done to trigger tc.Release
	globalLeaderElector.Lock()
	cancel()
	globalLeaderElector.Unlock()
	time.Sleep(time.Millisecond * 10)

	globalLeaderElector.Lock()
	assert.Equal(t, 1, tc.ReleaseCalledTimes)
	globalLeaderElector.Unlock()
}
