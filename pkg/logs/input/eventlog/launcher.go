// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package eventlog

import (
	"fmt"

	log "github.com/cihub/seelog"

	"github.com/DataDog/datadog-agent/pkg/logs/auditor"
	"github.com/DataDog/datadog-agent/pkg/logs/config"
	"github.com/DataDog/datadog-agent/pkg/logs/pipeline"
)

func generateIdentifier(channelPath, query string) string {
	return fmt.Sprintf("%s;%s", channelPath, query)
}

// Launcher is available only on windows
type Launcher struct {
	sources          []*config.LogSource
	pipelineProvider pipeline.Provider
	auditor          *auditor.Auditor
	tailers          map[string]*Tailer
}

// New returns a new Launcher.
func New(sources []*config.LogSource, pipelineProvider pipeline.Provider, auditor *auditor.Auditor) *Launcher {
	windowsEventSources := []*config.LogSource{}
	for _, source := range sources {
		if source.Config.Type == config.EventLogType {
			windowsEventSources = append(windowsEventSources, source)
		}
	}
	return &Launcher{
		sources:          windowsEventSources,
		pipelineProvider: pipelineProvider,
		auditor:          auditor,
		tailers:          make(map[string]*Tailer),
	}
}

// Start starts new tailers.
// FIXME: actually scan for channels?
func (l *Launcher) Start() {
	log.Info("Start tailing eventlog")

	for _, source := range l.sources {
		identifier := generateIdentifier(source.Config.ChannelPath, source.Config.Query)
		if _, exists := l.tailers[identifier]; exists {
			// tailer already setup
			continue
		}
		tailer, err := l.setupTailer(source)
		// FIXME: Agent status
		if err != nil {
			log.Warn("Could not set up tailer: ", err)
		} else {
			l.tailers[identifier] = tailer
		}
		tailer.Start()
	}
}

// Stop stops all active tailers
func (l *Launcher) Stop() {
	// stopper := restart.NewParallelStopper()
	// for _, tailer := range l.tailers {
	// 	stopper.Add(tailer)
	// 	delete(l.tailers, tailer.Identifier())
	// }
	// stopper.Stop()
}

// setupTailer configures and starts a new tailer,
// returns the tailer or an error.
func (l *Launcher) setupTailer(source *config.LogSource) (*Tailer, error) {
	query := source.Config.Query
	if query == "" {
		query = "*"
	}

	channelPath := source.Config.ChannelPath
	if channelPath == "" {
		channelPath = "System" // FIXME
	}

	tailer := NewTailer(source, channelPath, query, l.pipelineProvider.NextPipelineChan())
	tailer.Start() // FIXME restart at last bookmark
	return tailer, nil
}
