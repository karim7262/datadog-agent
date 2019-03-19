// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

//+build zlib

package jsonstream

import (
	"bytes"

	"github.com/json-iterator/go"

	"github.com/DataDog/datadog-agent/pkg/forwarder"
	"github.com/DataDog/datadog-agent/pkg/serializer/marshaler"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// PayloadBuilder is used to build payloads. PayloadBuilder allocates memory based
// on what was previously need to serialize payloads. Keep that in mind and
// use multiple PayloadBuilders for different sources.
type PayloadBuilder struct {
	inputSizeHint, outputSizeHint int
}

// NewPayloadBuilder creates a new PayloadBuilder with default values.
func NewPayloadBuilder() *PayloadBuilder {
	return &PayloadBuilder{
		inputSizeHint:  4096,
		outputSizeHint: 4096,
	}
}

var marshaller = jsoniter.Config{
	EscapeHTML:                    false,
	ObjectFieldMustBeSimpleString: true,
}.Froze()

// Build serializes a metadata payload and sends it to the forwarder
func (b *PayloadBuilder) Build(m marshaler.StreamJSONMarshaler) (forwarder.Payloads, error) {
	var payloads forwarder.Payloads
	var i int
	itemCount := m.Len()
	expvarsTotalCalls.Add(1)

	// Inner buffers for the compressor
	input := bytes.NewBuffer(make([]byte, 0, b.inputSizeHint))
	output := bytes.NewBuffer(make([]byte, 0, b.outputSizeHint))

	stream := jsoniter.NewStream(marshaller, nil, 1024*8)

	compressor, err := newCompressor(input, output, m.JSONHeader(), m.JSONFooter())
	if err != nil {
		return nil, err
	}

	for i < itemCount {
		stream.Reset(nil)
		err := m.WriteItem(stream, i)
		if err != nil {
			log.Warnf("error marshalling an item, skipping: %s", err)
			i++
			continue
		}

		serializedItem := stream.Buffer()

		switch compressor.addItem(serializedItem) {
		case errPayloadFull:
			// payload is full, we need to create a new one
			payload, err := compressor.close()
			if err != nil {
				return payloads, err
			}
			payloads = append(payloads, &payload)
			input.Reset()
			output.Reset()
			compressor, err = newCompressor(input, output, m.JSONHeader(), m.JSONFooter())
			if err != nil {
				return nil, err
			}
		case nil:
			// All good, continue to next item
			i++
			expvarsTotalItems.Add(1)
			continue
		default:
			// Unexpected error, drop the item
			i++
			log.Warnf("Dropping an item, %s: %s", m.DescribeItem(i), err)
			expvarsItemDrops.Add(1)
			continue
		}
	}

	// Close last payload
	payload, err := compressor.close()
	if err != nil {
		return payloads, err
	}
	payloads = append(payloads, &payload)

	b.inputSizeHint = input.Cap()
	b.outputSizeHint = output.Cap()

	return payloads, nil
}
