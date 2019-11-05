package dogstatsd

import (
	"bytes"
)

type messageType int

const (
	metricSampleType messageType = iota
	serviceCheckType
	eventType
)

var (
	eventPrefix        = []byte("_e{")
	serviceCheckPrefix = []byte("_sc")

	fieldSeparator = []byte("|")
	colonSeparator = []byte(":")
	commaSeparator = []byte(",")
)

const maxTags = 128

type parsedTags struct {
	tags      [maxTags][]byte
	tagsCount int
}

func findMessageType(message []byte) messageType {
	if bytes.HasPrefix(message, eventPrefix) {
		return eventType
	} else if bytes.HasPrefix(message, serviceCheckPrefix) {
		return serviceCheckType
	}
	// Note that random gibberish is interpreted as a metric since they don't
	// contain any easily identifiable feature
	return metricSampleType
}

// nextField returns the data found before the first fieldSeparator and
// the remainder, as a no-heap alternative to bytes.Split.
// If the separator is not found, the remainder is nil.
func nextField(message []byte) ([]byte, []byte) {
	return nextFieldSeparator(message, fieldSeparator)
}

func nextFieldSeparator(message, separator []byte) ([]byte, []byte) {
	sepIndex := bytes.Index(message, separator)
	if sepIndex == -1 {
		return message, nil
	}
	return message[:sepIndex], message[sepIndex+1:]
}

func parseTags(rawTags []byte) parsedTags {
	if len(rawTags) == 0 {
		return parsedTags{}
	}
	tags := parsedTags{}

	var tag []byte
	remainder := rawTags
	for tags.tagsCount < maxTags {
		tag, remainder = nextFieldSeparator(remainder, commaSeparator)
		tags.tags[tags.tagsCount] = tag
		tags.tagsCount++
		if remainder == nil {
			break
		}
	}
	return tags
}
