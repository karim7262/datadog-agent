// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2020 Datadog, Inc.

// +build clusterchecks

package cloudfoundry

import (
	"context"
	"sync"
	"time"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/locket"
	"code.cloudfoundry.org/locket/models"

	"github.com/DataDog/datadog-agent/pkg/clusteragent/clusterchecks/types"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	NotALeaderResponse    = "not me"
	LocketClusterAgentKey = "datadog-cluster-agent-leader"
)

type LeaderElector struct {
	sync.Mutex
	cancelContext   context.Context
	configured      bool
	hostname        string
	leaderIP        string
	locketAPIClient models.LocketClient
	refreshInterval time.Duration
}

var (
	globalLeaderElector     *LeaderElector = &LeaderElector{}
	globalLeaderElectorLock sync.Mutex
)

func RunGlobalLeaderElector(ctx context.Context, locketURL, cafile, certfile, keyfile, hostname string, refreshInterval time.Duration, testing models.LocketClient) (types.LeaderIPCallback, error) {
	globalLeaderElectorLock.Lock()
	defer globalLeaderElectorLock.Unlock()
	var err error

	if globalLeaderElector.configured {
		return globalLeaderElector.GetLeaderIP, nil
	}

	cfg := locket.ClientLocketConfig{
		LocketAddress:        locketURL,
		LocketCACertFile:     cafile,
		LocketClientCertFile: certfile,
		LocketClientKeyFile:  keyfile,
	}

	logger := lager.NewLogger("locket")
	if testing != nil {
		globalLeaderElector.locketAPIClient = testing
	} else {
		if globalLeaderElector.locketAPIClient, err = locket.NewClient(logger, cfg); err != nil {
			return nil, err
		}

	}
	globalLeaderElector.cancelContext = ctx
	globalLeaderElector.configured = true
	globalLeaderElector.hostname = hostname
	globalLeaderElector.leaderIP = NotALeaderResponse
	globalLeaderElector.refreshInterval = refreshInterval
	go globalLeaderElector.run()

	return globalLeaderElector.GetLeaderIP, nil
}

func (e *LeaderElector) GetLeaderIP() (types.LeaderIPResult, error) {
	e.Lock()
	defer e.Unlock()
	return types.LeaderIPResult{
		IP:           e.leaderIP,
		Redirectable: false,
	}, nil
}

func (e *LeaderElector) run() {
	// NOTE: we try to acquire every half of the refresh interval - this ensures that:
	// * the leader doesn't randomly switch - as long as it's running correctly, it will always reacquire before
	//   the TTL of lock expires (as TTL == refresh interval)
	// * when a leader fails/shuts down/disappears, a new leader acquires a lock reasonably fast
	dataRefreshTicker := time.NewTicker(e.refreshInterval / 2)
	e.tryAcquireLock()
	for {
		select {
		case <-dataRefreshTicker.C:
			e.tryAcquireLock()
		case <-e.cancelContext.Done():
			dataRefreshTicker.Stop()
			// when shutting down, release the lock explicitly so that other replica can grab it ASAP
			if e.leaderIP == "" {
				e.releaseLock()
			}
			return
		}
	}
}

func (e *LeaderElector) releaseLock() {
	log.Debug("Releasing leader lock ...")
	e.Lock()
	defer e.Unlock()
	e.leaderIP = NotALeaderResponse
	rr := &models.ReleaseRequest{
		Resource: &models.Resource{
			Key:   LocketClusterAgentKey,
			Owner: e.hostname,
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := e.locketAPIClient.Release(ctx, rr)
	if err != nil {
		log.Warnf("Failed to release leader lock on shutdown: %s", err.Error())
	} else {
		log.Info("Successfully released leader lock")
	}
}

func (e *LeaderElector) tryAcquireLock() {
	log.Debug("Trying to acquire leader lock ...")
	e.Lock()
	defer e.Unlock()
	ttl := int64(e.refreshInterval.Seconds())
	lr := &models.LockRequest{
		Resource: &models.Resource{
			Key:      LocketClusterAgentKey,
			Owner:    e.hostname,
			Value:    e.hostname,
			TypeCode: models.LOCK,
		},
		TtlInSeconds: ttl,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := e.locketAPIClient.Lock(ctx, lr)
	if err != nil {
		e.leaderIP = NotALeaderResponse
		log.Infof("Failed to acquire leader lock %s", err.Error())
	} else {
		e.leaderIP = ""
		log.Infof("Successfully acquired leader lock for %d seconds", ttl)
	}
}
