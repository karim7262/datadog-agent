/*
 * Unless explicitly stated otherwise all files in this repository are licensed
 * under the Apache License Version 2.0.
 * This product includes software developed at Datadog (https://www.datadoghq.com/).
 * Copyright 2016-2019 Datadog, Inc.
 */

package kubernetes

import (
	"github.com/DataDog/datadog-agent/pkg/logs/message"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPartialLineHandler_HandleFullMessage(t *testing.T) {
	h := preparedPartialLineHandler()
	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout F first full message"))
	var output *message.Message
	output = <- h.outputChan
	assert.Equal(t, "2019-05-15T13:34:26.011123506Z", output.Timestamp)
	assert.Equal(t, "first full message", string(output.Content))
	h.Stop()
}

func TestPartialLineHandler_HandleMultipleMessages(t *testing.T) {
	h := preparedPartialLineHandler()
	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout P first part message"))
	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout P second part message"))
	h.Handle([]byte("2019-05-15T13:34:26.011123507Z stdout F final part message"))
	var output *message.Message
	output = <- h.outputChan
	assert.Equal(t, "2019-05-15T13:34:26.011123507Z", output.Timestamp)
	assert.Equal(t, "first part messagesecond part messagefinal part message", string(output.Content))
	h.Stop()
}

func TestPartialLineHandler_TimeoutSendPartialMessage(t *testing.T) {
	h := preparedPartialLineHandler()
	var output *message.Message

	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout P first part message"))
	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout P second part message"))

	time.Sleep(2*time.Second)
	h.Handle([]byte("2019-05-15T13:34:26.011123507Z stdout F final part message"))


	output = <- h.outputChan
	assert.Equal(t, "2019-05-15T13:34:26.011123506Z", output.Timestamp)
	assert.Equal(t, "first part messagesecond part message", string(output.Content))

	output = <- h.outputChan
	assert.Equal(t, "2019-05-15T13:34:26.011123507Z", output.Timestamp)
	assert.Equal(t, "final part message", string(output.Content))
	h.Stop()
}

func TestPartialLineHandler_InvalidMessage(t *testing.T) {
	h := preparedPartialLineHandler()
	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout I unexpected message flag"))
	output := <- h.outputChan
	assert.Equal(t, "unexpected message flag", string(output.Content))
}

func TestPartialLineHandler_InvalidMessageShouldBeAlone(t *testing.T) {
	h := preparedPartialLineHandler()
	var output *message.Message

	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout P first part message"))
	h.Handle([]byte("2019-05-15T13:34:26.011123506Z stdout I unexpected message flag"))

	output = <- h.outputChan
	assert.Equal(t, "2019-05-15T13:34:26.011123506Z", output.Timestamp)
	assert.Equal(t, "first part message", string(output.Content))

	output = <- h.outputChan
	assert.Equal(t, "unexpected message flag", string(output.Content))
}

func preparedPartialLineHandler() *PartialLineHandler {
	outputChan := make(chan *message.Message, 10)
	var h = NewPartialLineHandler(outputChan, 10*time.Microsecond, kubernetes.Parser)
	h.Start()
	return h
}
