// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/cihub/seelog"
	yaml "gopkg.in/yaml.v2"
)

func parseFiles(names []string) ([]configFormat, error) {
	var configs []configFormat

	for _, file := range names {
		// Check name
		checkName, err := getCheckName(file)
		if err != nil {
			log.Infof("Skipping %s: %s", file, err)
			continue
		}
		// Read file contents
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			log.Warnf("Skipping %s: %s", file, err)
			continue
		}
		// Parse configuration
		var c configFormat
		err = yaml.Unmarshal(yamlFile, &c)
		if err != nil {
			log.Warnf("Skipping %s: %s", file, err)
			continue
		}

		c.CheckName = checkName
		configs = append(configs, c)
	}
	return configs, nil
}

func getCheckName(fileName string) (string, error) {
	fileBase := filepath.Base(fileName)
	ext := filepath.Ext(fileBase)
	if ext != ".yaml" && ext != ".yml" {
		return "", fmt.Errorf("invalid file extension: %q", ext)
	}
	checkName := strings.TrimSuffix(fileBase, ext)
	if checkName == "" {
		return "", fmt.Errorf("invalid file name: %q", fileBase)
	}

	return strings.TrimSuffix(fileBase, ext), nil
}
