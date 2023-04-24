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

import "testing"

func Test_calculateMakeConcurrency(t *testing.T) {
	type testCase struct {
		input              string
		expectedMakeJobs   string
		expectedCPURequest string
	}

	for _, test := range []testCase{
		{
			input: "3500m",

			expectedMakeJobs:   "-j3",
			expectedCPURequest: "3500m",
		},
		{
			input: "5500M",

			expectedMakeJobs:   "-j5",
			expectedCPURequest: "5500m",
		},
		{
			input: "55",

			expectedMakeJobs:   "-j55",
			expectedCPURequest: "55",
		},
		{
			input: "55000m",

			expectedMakeJobs:   "-j55",
			expectedCPURequest: "55000m",
		},
		{
			input: "500m",

			expectedMakeJobs:   "-j1",
			expectedCPURequest: "500m",
		},
		{
			input: "0.5",

			expectedMakeJobs:   "-j1",
			expectedCPURequest: "0.5",
		},
	} {
		gotMakeJobs, gotCPURequest := calculateMakeConcurrency(test.input)

		if gotMakeJobs != test.expectedMakeJobs {
			t.Errorf("make --jobs: expected %q but got %q", test.expectedMakeJobs, gotMakeJobs)
		}

		if gotCPURequest != test.expectedCPURequest {
			t.Errorf("CPU request: expected %q but got %q", test.expectedCPURequest, gotCPURequest)
		}
	}
}

func Test_calculateMakeConcurrency_NoRequest_Failure(t *testing.T) {
	caughtPanic := false

	defer func() {
		if r := recover(); r != nil {
			caughtPanic = true
		}

		if !caughtPanic {
			t.Fatalf("expected a panic for with no CPU request for calculateMakeConcurrency but didn't get one")
		}
	}()

	calculateMakeConcurrency("")
}

func Test_calculateMakeConcurrency_InvalidRequestMillis_Failure(t *testing.T) {
	caughtPanic := false

	defer func() {
		if r := recover(); r != nil {
			caughtPanic = true
		}

		if !caughtPanic {
			t.Fatalf("expected a panic for with no CPU request for calculateMakeConcurrency but didn't get one")
		}
	}()

	calculateMakeConcurrency("100am")
}

func Test_calculateMakeConcurrency_InvalidRequestCPUs_Failure(t *testing.T) {
	caughtPanic := false

	defer func() {
		if r := recover(); r != nil {
			caughtPanic = true
		}

		if !caughtPanic {
			t.Fatalf("expected a panic for with no CPU request for calculateMakeConcurrency but didn't get one")
		}
	}()

	calculateMakeConcurrency("1a")
}

func Test_splitKubernetesVersion(t *testing.T) {
	type testCase struct {
		input string

		expectedMajorVersion int
		expectedMinorVersion int

		expectError bool
	}

	for _, test := range []testCase{
		{
			input: "1.23",

			expectedMajorVersion: 1,
			expectedMinorVersion: 23,

			expectError: false,
		},
		{
			input: "1.23.1",

			expectedMajorVersion: 1,
			expectedMinorVersion: 23,

			expectError: false,
		},
		{
			input: "123",

			expectedMajorVersion: 0,
			expectedMinorVersion: 0,

			expectError: true,
		},
		{
			input: "2.24",

			expectedMajorVersion: 2,
			expectedMinorVersion: 24,

			expectError: false,
		},
		{
			input: "1a.24",

			expectedMajorVersion: 0,
			expectedMinorVersion: 0,

			expectError: true,
		},
		{
			input: "1.a24",

			expectedMajorVersion: 0,
			expectedMinorVersion: 0,

			expectError: true,
		},
		{
			input: "a1.a24",

			expectedMajorVersion: 0,
			expectedMinorVersion: 0,

			expectError: true,
		},
	} {
		gotMajor, gotMinor, err := splitKubernetesVersion(test.input)

		if (err != nil) != test.expectError {
			t.Errorf("expectError=%v, err=%v", test.expectError, err)
		}

		if gotMajor != test.expectedMajorVersion {
			t.Errorf("got major version %q from %q, wanted %q", gotMajor, test.input, test.expectedMajorVersion)
		}

		if gotMinor != test.expectedMinorVersion {
			t.Errorf("got minor version %q from %q, wanted %q", gotMinor, test.input, test.expectedMinorVersion)
		}
	}
}
