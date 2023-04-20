// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/twitchdev/twitch-cli/cmd"
	"github.com/twitchdev/twitch-cli/internal/util"
)

var buildVersion string

const (
	CPU_PROFILER_BOOL_ENV_VARIABLE     = "TWITCH_CLI_ENABLE_CPU_PROFILER"
	CPU_PROFILER_FILENAME_ENV_VARIABLE = "TWITCH_CLI_CPU_PROFILER_FILE"
	CPU_PROFILER_DEFAULT_FILENAME      = "cpu.prof"
)

func main() {
	enableCpuProfiler, err := strconv.ParseBool(os.Getenv(CPU_PROFILER_BOOL_ENV_VARIABLE))
	if err != nil {
		enableCpuProfiler = false
	}

	// Enable CPU profiler
	if enableCpuProfiler {
		cpuProfilerFilename := os.Getenv(CPU_PROFILER_FILENAME_ENV_VARIABLE)
		if strings.TrimSpace(cpuProfilerFilename) == "" {
			cpuProfilerFilename = CPU_PROFILER_DEFAULT_FILENAME
		}

		f, err := os.Create(cpuProfilerFilename)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer f.Close()

		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}

		defer pprof.StopCPUProfile()
	}

	if len(buildVersion) > 0 {
		util.SetVersion(buildVersion)
	}
	cmd.Execute()
}
