// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build windows

package eventlog

import (
	"C"
	"fmt"

	log "github.com/cihub/seelog"
)

// Identifier returns the unique identifier of the current event log stream being tailed.
func (t *Tailer) Identifier() string {
	return fmt.Sprintf("eventlog:%s", generateIdentifier(channelPath, query))
}

// Start starts tailing the event log
func (t *Tailer) Start() error {
	log.Info("Start tailing eventlog for channel %s query %s", channelPath, query)
	// go t.tail()
	return nil
}

func (t *Tailer) tail() {
	C.startEventSubscribe(
		C.CString(t.channelPath),
		C.CString(t.query),
		C.ULONGLONG(0),
		C.int(EvtSubscribeToFutureEvents),
	)
	// <-t.stop
	// log.Info("Stop tailing eventlog for channel %s query %s", channelPath, query)
	// t.done <- struct{}{}
}
