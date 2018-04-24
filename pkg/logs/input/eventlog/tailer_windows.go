// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

// +build windows

package eventlog

import (
	"C"

	log "github.com/cihub/seelog"
)
import "time"

// tail waits for message stop
func (t *Tailer) Start() {
	log.Warn("Starting event log tailing for channel ", t.config.ChannelPath, "query ", t.config.Query)
	go t.tail()
}

// Stop stops the tailer
func (t *Tailer) Stop() {
	log.Warn("Stop tailing event log")
	t.stop <- struct{}{}
	<-t.done
}

// tail waits for message stop
func (t *Tailer) tail() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-t.stop:
			t.done <- struct{}{}
			return
		case <-ticker.C:
			log.Warn("tailing file")
		}
	}
}

func (t *Tailer) tail2() {
	C.startEventSubscribe(
		C.CString(t.channelPath),
		C.CString(t.query),
		C.ULONGLONG(0),
		C.int(EvtSubscribeToFutureEvents),
	)
}
