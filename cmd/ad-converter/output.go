// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package main

import (
	"encoding/json"
	"fmt"
)

func (c configRawMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}(c))
}

func jsonizeConfig(configs []configFormat) (jsonConfigs, error) {
	var names []string
	var initConfigs []configRawMap
	var instances []configRawMap

	for _, c := range configs {
		if c.InitConfig == nil {
			c.InitConfig = make(configRawMap)
		}
		for _, i := range c.Instances {
			if i == nil {
				i = make(configRawMap)
			}
			names = append(names, c.CheckName)
			initConfigs = append(initConfigs, c.InitConfig)
			instances = append(instances, i)
		}
	}

	var j jsonConfigs
	b, err := json.Marshal(names)
	if err != nil {
		return j, err
	}
	j.Names = string(b)
	b, err = json.Marshal(initConfigs)
	if err != nil {
		return j, err
	}
	j.InitConfigs = string(b)
	b, err = json.Marshal(instances)
	if err != nil {
		return j, err
	}
	j.Instances = string(b)

	return j, nil
}

func printDockerLabels(conf jsonConfigs) {
	fmt.Printf("LABEL \"com.datadoghq.ad.check_names\"='%s'\n", conf.Names)
	fmt.Printf("LABEL \"com.datadoghq.ad.init_configs\"='%s'\n", conf.InitConfigs)
	fmt.Printf("LABEL \"com.datadoghq.ad.instances\"='%s'\n", conf.Instances)
}

func printComposeLabels(conf jsonConfigs) {
	fmt.Printf("com.datadoghq.ad.check_names: '%s'\n", conf.Names)
	fmt.Printf("com.datadoghq.ad.init_configs: '%s'\n", conf.InitConfigs)
	fmt.Printf("com.datadoghq.ad.instances: '%s'\n", conf.Instances)
}

func printKubeAnnotations(conf jsonConfigs, containerName string) {
	fmt.Printf("ad.datadoghq.com/%s.check_names: '%s'\n", containerName, conf.Names)
	fmt.Printf("ad.datadoghq.com/%s.init_configs: '%s'\n", containerName, conf.InitConfigs)
	fmt.Printf("ad.datadoghq.com/%s.instances: '%s'\n", containerName, conf.Instances)
}
