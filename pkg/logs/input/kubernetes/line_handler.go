/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache License Version 2.0.
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2016-2019 Datadog, Inc.
 */

package kubernetes

import (
	"bytes"
	"fmt"
	"github.com/DataDog/datadog-agent/pkg/logs/message"
	"github.com/DataDog/datadog-agent/pkg/logs/decoder"
	"time"
)

var (
	fullMessage = []byte{'F'}
	partialMessage = []byte{'P'}
)

type PartialLineHandler struct {
	lineChan chan []byte
	outputChan chan *message.Message
	lineBuffer *decoder.LineBuffer
	flushTimeout time.Duration
	psr *parser
	lastSeenTimestamp string
}

func NewPartialLineHandler(outputChan chan *message.Message, flushTimeout time.Duration) *PartialLineHandler {
	return &PartialLineHandler{
		lineChan: make(chan []byte),
		outputChan: outputChan,
		lineBuffer: decoder.NewLineBuffer(),
		flushTimeout: flushTimeout,
		psr: Parser,
	}
}

// Handle forward lines to lineChan to process them
func (h *PartialLineHandler) Handle(content []byte) {
	h.lineChan <- content
}

// Stop the lineHandler from processing lines
func (h *PartialLineHandler) Stop() {
	close(h.lineChan)
}

// Start the handler
func (h *PartialLineHandler) Start() {
	go h.run()
}

func (h *PartialLineHandler) run() {
	// start timer
	flushTimer := time.NewTimer(h.flushTimeout)
	// schedule a clean up
	defer func() {
		flushTimer.Stop()
		close(h.outputChan)
	}()

	for {
		select {
		case line, isOpen := <-h.lineChan:
			if !isOpen { // lineChan closed.
				return
			}
			if !flushTimer.Stop() {
				select {
				case <-flushTimer.C:
				default:
				}
			}
			h.process(line)
			flushTimer.Reset(h.flushTimeout)
		case <-flushTimer.C: // flush time out
			h.sendContent()
			//flushTimer.Reset(h.flushTimeout)
		}
	}
}

func (h *PartialLineHandler) process(line []byte) {
	content, status, timestamp, flag, err = h.psr.Split(line)
	h.lastSeenTimestamp = timestamp
	h.feedBuffer(line, content)
	if bytes.Equal(flag, fullMessage) {
		h.sendContent()
	} else if !bytes.Equal(flag, partialMessage) {
		log.Debug(fmt.Errorf("unrecognized flag: %v", string(flag)))
		if !h.lineBuffer.IsEmpty() {
			h.sendContent()
		}
		h.lineBuffer.AddIncompleteLine(line)
		h.sendContent()
	}
}

func (h *PartialLineHandler) feedBuffer(line []byte, msg []byte) {
	if h.lineBuffer.IsEmpty() {
		h.lineBuffer.AddIncompleteLine(line)
	} else {
		h.lineBuffer.AddIncompleteLine(msg)
	}
}

func (h *PartialLineHandler) sendContent() {
	defer h.lineBuffer.Reset()
	content, rawDatalen := h.lineBuffer.Content()
	content = bytes.TrimSpace(content)
	if len(content) > 0 {
		output, err := h.psr.Parse(content)
		if err != nil {
			log.Debug(err)
		}
		if output != nil && len(output.Content) > 0 {
			output.Timestamp = h.lastSeenTimestamp
			output.RawDataLen = rawDatalen
			h.outputChan <- output
		}
	}
}
