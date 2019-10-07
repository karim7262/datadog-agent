// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package telemetry

// Counter TODO(remy): doc
type Counter interface {
	// TODO(remy): doc
	Inc(tags ...string)
	// TODO(remy): doc
	Add(value float64, tags ...string)
}
