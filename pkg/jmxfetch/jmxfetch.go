// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

// +build jmx

package jmxfetch

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	api "github.com/DataDog/datadog-agent/pkg/api/util"
	"github.com/DataDog/datadog-agent/pkg/status/health"
)

/*
#cgo CFLAGS: -I/usr/lib/jvm/java-13-openjdk-amd64/include/ -I/usr/lib/jvm/java-13-openjdk-amd64/include/linux
#cgo LDFLAGS: -L/usr/lib/jvm/java-13-openjdk-amd64/lib/server/ -ljvm

#include "java_class_loader.h"
*/
import "C"

const (
	jmxJarName                        = "jmxfetch.jar"
	jmxMainClass                      = "org.datadog.jmxfetch.App"
	defaultJmxCommand                 = "collect"
	defaultJvmMaxMemoryAllocation     = " -Xmx200m"
	defaultJvmInitialMemoryAllocation = " -Xms50m"
	jvmCgroupMemoryAwareness          = " -XX:+UnlockExperimentalVMOptions -XX:+UseCGroupMemoryLimitForHeap"
	defaultJavaBinPath                = "java"
	defaultLogLevel                   = "info"
)

var (
	jmxLogLevelMap = map[string]string{
		"trace":    "TRACE",
		"debug":    "DEBUG",
		"info":     "INFO",
		"warn":     "WARN",
		"warning":  "WARN",
		"error":    "ERROR",
		"err":      "ERROR",
		"critical": "FATAL",
	}
	jvmCgroupMemoryIncompatOptions = []string{
		"Xmx",
		"XX:MaxHeapSize",
		"Xms",
		"XX:InitialHeapSize",
	}
)

// JMXFetch represent a jmxfetch instance.
type JMXFetch struct {
	JavaBinPath        string
	JavaOptions        string
	JavaToolsJarPath   string
	JavaCustomJarPaths []string
	LogLevel           string
	Command            string
	ReportOnConsole    bool
	Checks             []string
	IPCPort            int
	IPCHost            string
	defaultJmxCommand  string
	cmd                *exec.Cmd
	managed            bool
	shutdown           chan struct{}
	stopped            chan struct{}
}

func (j *JMXFetch) setDefaults() {
	if j.JavaBinPath == "" {
		j.JavaBinPath = defaultJavaBinPath
	}
	if j.JavaCustomJarPaths == nil {
		j.JavaCustomJarPaths = []string{}
	}
	if j.LogLevel == "" {
		j.LogLevel = defaultLogLevel
	}
	if j.Command == "" {
		j.Command = defaultJmxCommand
	}
	if j.Checks == nil {
		j.Checks = []string{}
	}
}

// Start starts the JMXFetch process
func (j *JMXFetch) Start(manage bool) error {
	csubprocessArgs := []string{"-Djava.class.path=.:/home/maxime/dev/datadog-agent/pkg/jmxfetch/:/home/maxime/dev/jmxfetch/target/jmxfetch-0.30.0-jar-with-dependencies.jar"}
	os.Setenv("SESSION_TOKEN", api.GetAuthToken())
	ac := C.int(len(csubprocessArgs))
	av := make([]*C.char, ac);
	for i := range csubprocessArgs {
		av[i] = C.CString(csubprocessArgs[i])
	}

	go C.run_jmx(ac, &av[0])

	return nil;
}

// Wait waits for the end of the JMXFetch process and returns the error code
func (j *JMXFetch) Wait() error {
	return j.cmd.Wait()
}

func (j *JMXFetch) heartbeat(beat *time.Ticker) {
	health := health.Register("jmxfetch")
	defer health.Deregister()

	for range beat.C {
		select {
		case <-health.C:
		case <-j.shutdown:
			return
		}
	}
}

// Up returns if JMXFetch is up - used by healthcheck
func (j *JMXFetch) Up() (bool, error) {
	// TODO: write windows implementation
	process, err := os.FindProcess(j.cmd.Process.Pid)
	if err != nil {
		return false, fmt.Errorf("Failed to find process: %s\n", err)
	}

	// from man kill(2):
	// if sig is 0, then no signal is sent, but error checking is still performed
	err = process.Signal(syscall.Signal(0))
	return err == nil, err
}
