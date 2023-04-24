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
	"testing"
)

func Test_addStandardE2ELabels_NewKubernetes(t *testing.T) {
	// on any version of k8s greater than or equal to 1.23 we should enable all feature gates
	// and use SSA
	for _, testVersion := range []string{"1.22", "1.23", "1.24", "2.1"} {
		test := jobTemplate(
			"test-test",
			"some description",
			addStandardE2ELabels(testVersion),
		)

		if borkedValue, ok := test.Labels["preset-enable-all-feature-gates-disable-ssa"]; ok {
			t.Errorf("didn't expect 'preset-enable-all-feature-gates-disable-ssa' to be set for newer k8s version %q but it has value %s", testVersion, borkedValue)
		}

		gatesLabel, ok := test.Labels["preset-enable-all-feature-gates"]
		if !ok {
			t.Errorf("missing 'preset-enable-all-feature-gates' label after addStandardE2ELabels for newer k8s %s", testVersion)
			continue
		}

		if gatesLabel != "true" {
			t.Errorf("expected feature gates label to be 'true' but it was %q", gatesLabel)
		}
	}
}

func Test_addStandardE2ELabels_OldKubernetes(t *testing.T) {
	// on any version of k8s greater than or equal to 1.23 we should enable all feature gates
	// but disable SSA
	for _, testVersion := range []string{"1.21", "1.20", "1.0"} {
		test := jobTemplate(
			"test-test",
			"some description",
			addStandardE2ELabels(testVersion),
		)

		if borkedValue, ok := test.Labels["preset-enable-all-feature-gates"]; ok {
			t.Errorf("didn't expect 'preset-enable-all-feature-gates' to be set for older k8s version %q but it has value %s", testVersion, borkedValue)
		}

		gatesLabel, ok := test.Labels["preset-enable-all-feature-gates-disable-ssa"]
		if !ok {
			t.Errorf("missing 'preset-enable-all-feature-gates-disable-ssa' label after addStandardE2ELabels for older k8s version %s", testVersion)
			continue
		}

		if gatesLabel != "true" {
			t.Errorf("expected feature gates label to be 'true' but it was %q", gatesLabel)
		}
	}
}

func Test_addStandardE2ELabels_ProgrammerError(t *testing.T) {
	k8sVersion := "1a.2a3a4a"
	caughtPanic := false

	defer func() {
		if r := recover(); r != nil {
			caughtPanic = true
		}

		if !caughtPanic {
			t.Fatalf("expected a panic for addStandardE2ELabels with k8s version %q but didn't get one", k8sVersion)
		}
	}()

	// programmer error with k8s version should panic
	jobTemplate(
		"test-test",
		"some description",
		addStandardE2ELabels(k8sVersion),
	)
}
