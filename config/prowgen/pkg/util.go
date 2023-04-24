// +skip_license_check
/*
Copyright 2022 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	milliCPUToCPU = 1000.0
)

func calculateMakeConcurrency(cpuRequest string) (string, string) {
	if len(cpuRequest) == 0 {
		panic("cannot determine value for NUM in make -j<NUM> without a configured CPU request")
	}

	cpuMultiplier := milliCPUToCPU

	cpuRequest = strings.ToLower(cpuRequest)

	originalCPURequest := cpuRequest

	if strings.HasSuffix(cpuRequest, "m") {
		cpuRequest = strings.TrimSuffix(cpuRequest, "m")
		cpuMultiplier = 1.0
	}

	parsedCPUs, err := strconv.ParseFloat(cpuRequest, 64)
	if err != nil {
		panic(fmt.Errorf("CPU request %q wasn't a number: %w", originalCPURequest, err))
	}

	milliCPUs := parsedCPUs * cpuMultiplier

	makeJobs := int(math.Floor(milliCPUs / milliCPUToCPU))

	if makeJobs < 1 {
		makeJobs = 1
	}

	return fmt.Sprintf("-j%d", makeJobs), originalCPURequest
}

func splitKubernetesVersion(version string) (int, int, error) {
	versionParts := strings.Split(version, ".")
	if len(versionParts) == 1 {
		return 0, 0, fmt.Errorf("invalid version format %q; wanted at least two parts separated by a '.'", version)
	}

	majorPart, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid major version %q: %w", versionParts[0], err)
	}

	minorPart, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid minor version %q: %w", versionParts[1], err)
	}

	return majorPart, minorPart, nil
}
