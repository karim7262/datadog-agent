// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package telemetry

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-agent/pkg/config"
)

func dogstatsdUDPDropCount() (int, error) {
	f, err := os.Open("/proc/net/udp")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return readDogstatsdUDPDrops(f)
}

func readDogstatsdUDPDrops(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)

	// pop header
	if !scanner.Scan() {
		return 0, errors.New("/proc/net/udp is empty")
	}

	const localAddressIndex = 1
	const dropsIndex = 12

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)

		rawLocalAddress := fields[localAddressIndex]

		port, err := strconv.ParseInt(rawLocalAddress[len(rawLocalAddress)-4:], 16, 32)
		if err != nil {
			return 0, errors.New("/proc/net/udp is not formatted correctly")
		}

		if port == config.Datadog.GetInt64("dogstatsd_port") {
			drops, err := strconv.ParseInt(fields[dropsIndex], 10, 32)
			return int(drops), err
		}
	}
	return 0, errors.New("dogstatsd socket was not found")
}
