// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package config

import (
	"flag"
	"path/filepath"
)

// Some dependencies (including the Kubernetes libraries)
// use the glog logging package, which does not allow to
// programatically configure it.

// This file hacks the program flags to configure the following
// glog behaviour:
//   - disable stderr output
//   - output logs to file, in the same folder as the agent
//   - default level 0, configurable via glog_verbosity_level

func setupGlog(logFile string) {
	// Avoid logging glog on stderr
	flag.Set("logtostderr", "false")

	// Set the verbosity for files. viper will cast the int to string for us
	flag.Set("v", Datadog.GetString("glog_verbosity_level"))

	// Output logs in the same folder as the agent (to get them in flare)
	flag.Set("log_dir", filepath.Dir(logFile))

	// Convinces goflags that we have called Parse() to avoid noisy logs.
	// OSS Issue: kubernetes/kubernetes#17162.
	flag.CommandLine.Parse([]string{})
}
