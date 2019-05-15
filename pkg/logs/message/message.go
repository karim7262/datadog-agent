// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package message

type Message interface {
	Content() []byte
	Origin() *Origin
	Status() string
	Timestamp() string
	RawDataLen() int
}

// Message represents a log line sent to datadog, with its metadata
type DefaultMessage struct {
	content    []byte
	origin     *Origin
	status     string
	timestamp  string
	rawDataLen int
}

// NewMessage returns a new message
func NewDefaultMessage(content []byte, origin *Origin, status string) *DefaultMessage {
	return &DefaultMessage{
		content: content,
		origin:  origin,
		status:  status,
	}
}

func (m *DefaultMessage) Status() string{
	if m.status == "" {
		m.status = StatusInfo
	}
	return m.status
}

// SetStatus sets the status of the message
func (m *DefaultMessage) SetStatus(status string) {
	m.status = status
}
