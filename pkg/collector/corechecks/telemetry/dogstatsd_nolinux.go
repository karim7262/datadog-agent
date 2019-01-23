// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build !linux

package telemetry

import "errors"

func dogstatsdDropCount() (int, error) {
	return 0, errors.New("Not implemented")
}
