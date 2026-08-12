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

package pkg

import (
	"testing"
)

// testProwContext returns a minimal ProwContext suitable for exercising
// presubmit/periodic registration.
func testProwContext() *ProwContext {
	return &ProwContext{
		Branch: "master",
		Org:    "cert-manager",
		Repo:   "cert-manager",
	}
}

// jobWithForbiddenLabels returns a job carrying every credential preset that
// presubmits must never mount, so we can assert they are rejected/retained.
func jobWithForbiddenLabels(name string) *Job {
	job := jobTemplate(name, "some description")
	for _, label := range presubmitForbiddenLabels {
		job.Labels[label] = "true"
	}
	return job
}

// Presubmits run unreviewed PR code, so addPresubmit must fail generation for
// a job carrying a credential preset, regardless of which generator added it.
func Test_addPresubmit_rejectsCredentialLabels(t *testing.T) {
	for _, label := range presubmitForbiddenLabels {
		t.Run(label, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("expected a panic for presubmit with credential preset %q, but generation succeeded", label)
				}
			}()

			job := jobTemplate("e2e", "some description")
			job.Labels[label] = "true"

			testProwContext().RequiredPresubmit(job)
		})
	}
}

// Periodics run merged, reviewed code, so they must keep their credentials:
// the rejection is a presubmit-only guard.
func Test_Periodics_retainsCredentialLabels(t *testing.T) {
	pc := testProwContext()
	pc.Periodics(jobWithForbiddenLabels("e2e"), 2)

	if len(pc.periodics) != 1 {
		t.Fatalf("expected exactly one periodic, got %d", len(pc.periodics))
	}

	labels := pc.periodics[0].Labels
	for _, label := range presubmitForbiddenLabels {
		if _, ok := labels[label]; !ok {
			t.Errorf("periodic must retain credential preset %q, but it was stripped", label)
		}
	}
}
