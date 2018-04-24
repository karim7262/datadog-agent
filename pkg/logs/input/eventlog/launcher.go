// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package eventlog

import (
	log "github.com/cihub/seelog"

	"github.com/DataDog/datadog-agent/pkg/logs/auditor"
	"github.com/DataDog/datadog-agent/pkg/logs/config"
	"github.com/DataDog/datadog-agent/pkg/logs/pipeline"
)

// Launcher is available only on windows
type Launcher struct {
	sources []*config.LogSource
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
		sources: windowsEventSources,
	}
}

// Start does nothing
func (l *Launcher) Start() {
	if len(l.sources) > 0 {
		log.Warn("WindowsEvent is not supported on this system.")
	}
}

// Stop does nothing
func (l *Launcher) Stop() {}
