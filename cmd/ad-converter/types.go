// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package main

type configRawMap map[string]interface{}

type configFormat struct {
	CheckName  string
	InitConfig configRawMap `yaml:"init_config"`
	LogsConfig configRawMap
	Instances  []configRawMap
}

type jsonConfigs struct {
	Names       string
	InitConfigs string
	Instances   string
}
