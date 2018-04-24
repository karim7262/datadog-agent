// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package eventlog

import (
	log "github.com/cihub/seelog"

	"github.com/DataDog/datadog-agent/pkg/logs/config"
	"github.com/DataDog/datadog-agent/pkg/logs/message"
)

// Tailer collects logs from a journal.
type Tailer struct {
	source      *config.LogSource
	channelPath string
	query       string
	outputChan  chan message.Message
	stop        chan struct{}
	done        chan struct{}
}

// NewTailer returns a new tailer.
func NewTailer(source *config.LogSource, channelPath string, query string, outputChan chan message.Message) *Tailer {
	return &Tailer{
		source:      source,
		channelPath: channelPath,
		query:       query,
		outputChan:  outputChan,
		stop:        make(chan struct{}, 1),
		done:        make(chan struct{}, 1),
	}
}

// setup does nothing
func (t *Tailer) setup() error {
	log.Info("EventLog is not supported on this system.")
	return nil
}

// tail waits for message stop
func (t *Tailer) Start() {
	<-t.stop
	t.done <- struct{}{}
}
