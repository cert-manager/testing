// +skip_license_check
/*
Copyright 2026 The cert-manager Authors.

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

package prowspecs

import (
	"strings"
	"testing"

	"prowgen/pkg"
)

// Jobs without a memory limit, and why. Every other generated job must set a
// memory limit equal to its memory request, so that a job which outgrows its
// request is killed in its own cgroup rather than taking down the node.
// See https://github.com/cert-manager/testing/issues/1240
var memoryLimitExemptions = map[string]string{
	"make-verify": "peak RSS is being re-measured with GOMAXPROCS capped; see issue 1240",
}

func Test_memoryLimitsMatchRequests(t *testing.T) {
	for _, branch := range KnownBranches() {
		spec, err := SpecForBranch(branch)
		if err != nil {
			t.Fatal(err)
		}
		jobFile := spec.GenerateJobFile()

		var jobs []*pkg.Job
		for _, presubmits := range jobFile.Presubmits {
			for _, p := range presubmits {
				jobs = append(jobs, &p.Job)
			}
		}
		for _, p := range jobFile.Periodics {
			jobs = append(jobs, &p.Job)
		}

		for _, job := range jobs {
			for _, c := range job.Spec.Containers {
				if exempt(job.Name) {
					if c.Resources.Limits != nil {
						t.Errorf("%s: %s is listed in memoryLimitExemptions but sets a memory limit; remove the exemption", branch, job.Name)
					}
					continue
				}
				if c.Resources.Limits == nil {
					t.Errorf("%s: %s has no memory limit; use memoryResources or add a reason to memoryLimitExemptions", branch, job.Name)
					continue
				}
				if c.Resources.Limits.Memory != c.Resources.Requests.Memory {
					t.Errorf("%s: %s memory limit %q does not equal request %q", branch, job.Name, c.Resources.Limits.Memory, c.Resources.Requests.Memory)
				}
			}
		}
	}
}

func exempt(jobName string) bool {
	for suffix := range memoryLimitExemptions {
		if strings.HasSuffix(jobName, suffix) {
			return true
		}
	}
	return false
}
