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

package prowspecs

import (
	"fmt"
	"slices"
	"strings"

	"prowgen/pkg"
)

// knownBranches specifies a BranchSpec for each possible branch to test against
// THIS IS WHAT YOU'RE MOST LIKELY TO NEED TO EDIT
// The branches and kubernetes versions below are likely to need to be updated after each cert-manager release!

// NB: There's at least one configurer (pkg/configurers.go) which will changes its operations
// based on the k8s version it's being run against.

var knownBranches map[string]BranchSpec = map[string]BranchSpec{
	"release-1.12": {
		prowContext: &pkg.ProwContext{
			Branch: "release-1.12",

			// Use latest image.
			Image: pkg.CommonTestImage,

			// NB: we don't use a presubmit dashboard outside of "master", currently
			PresubmitDashboard: false,
			PeriodicDashboard:  true,

			Org:  "cert-manager",
			Repo: "cert-manager",
		},

		primaryKubernetesVersion: "1.27",
		// NB: It would be nice to test 1.30 and 1.31 (and newer) here but newer versions of Kind don't
		// build images to support testing older k8s versions. E.g. kind v0.24.0 doesn't have images for
		// Kubernetes 1.24 and below
		otherKubernetesVersions: []string{"1.22", "1.23", "1.24", "1.25", "1.26", "1.28", "1.29"},

		e2eCPURequest:    "7000m",
		e2eMemoryRequest: "6Gi",

		// This older cert-manager release uses the ctl image to run the statupapicheck test
		containerNames: []string{"controller", "acmesolver", "ctl", "cainjector", "webhook"},

		// Keep using the old tests (for backwards compatibility)
		isPreMakefileModules: true,
	},
	"release-1.14": {
		prowContext: &pkg.ProwContext{
			Branch: "release-1.14",

			// Use latest image.
			Image: pkg.CommonTestImage,

			// NB: we don't use a presubmit dashboard outside of "master", currently
			PresubmitDashboard: false,
			PeriodicDashboard:  true,

			Org:  "cert-manager",
			Repo: "cert-manager",
		},

		primaryKubernetesVersion: "1.29",
		otherKubernetesVersions:  []string{"1.24", "1.25", "1.26", "1.27", "1.28"},

		e2eCPURequest:    "7000m",
		e2eMemoryRequest: "6Gi",

		// This older cert-manager release uses the NEW startupapicheck image to run the statupapicheck test
		// The release however still includes a ctl image (which is not used in the Helm chart)
		containerNames: []string{"controller", "acmesolver", "ctl", "startupapicheck", "cainjector", "webhook"},

		// Keep using the old tests (for backwards compatibility)
		isPreMakefileModules: true,
	},
	"release-1.15": {
		prowContext: &pkg.ProwContext{
			Branch: "release-1.15",

			// Use latest image.
			Image: pkg.CommonTestImage,

			// NB: we don't use a presubmit dashboard outside of "master", currently
			PresubmitDashboard: false,
			PeriodicDashboard:  true,

			Org:  "cert-manager",
			Repo: "cert-manager",
		},

		primaryKubernetesVersion: "1.30",

		// TODO: test k8s 1.31 here when possible; requires support in the release-1.15 branch on cert-manager
		otherKubernetesVersions: []string{"1.25", "1.26", "1.27", "1.28", "1.29", "1.31"},

		e2eCPURequest:    "7000m",
		e2eMemoryRequest: "6Gi",
	},
	"master": {
		prowContext: &pkg.ProwContext{
			Branch: "master",

			// Use latest image.
			Image: pkg.CommonTestImage,

			PresubmitDashboard: true,
			PeriodicDashboard:  true,

			Org:  "cert-manager",
			Repo: "cert-manager",
		},

		primaryKubernetesVersion: "1.31",
		otherKubernetesVersions:  []string{"1.27", "1.28", "1.29", "1.30"},

		e2eCPURequest:    "7000m",
		e2eMemoryRequest: "6Gi",
	},
}

// BranchSpec holds a specification of an entire test suite for a given branch, such as "master" or "release-1.9"
// That includes:
// - a ProwContext specifying things like the the repo, branch, dashboard names
// - the primary Kubernetes version (which is the version whose tests are always run for presubmits, among other uses)
// - the secondary Kubernetes versions, which are the rest of the supported versions for which tests should be generated
type BranchSpec struct {
	prowContext *pkg.ProwContext

	primaryKubernetesVersion string
	otherKubernetesVersions  []string

	e2eCPURequest    string
	e2eMemoryRequest string

	// TODO: remove this field once we've migrated to the new set of container names
	containerNames []string

	// TODO: remove this field once all versions use Makefile modules
	isPreMakefileModules bool
}

// GenerateJobFile will create a complete test file based on the BranchSpec. This
// assumes that all tests for all branches should be the same.
func (m *BranchSpec) GenerateJobFile() *pkg.JobFile {
	if !m.isPreMakefileModules {
		m.prowContext.RequiredPresubmit(pkg.MakeVerify(m.prowContext))
		m.prowContext.RequiredPresubmit(pkg.MakeTest(m.prowContext))
	} else {
		m.prowContext.RequiredPresubmit(pkg.MakeTestOld(m.prowContext))
		m.prowContext.RequiredPresubmit(pkg.ChartTestOld(m.prowContext))
	}

	for _, secondaryVersion := range m.otherKubernetesVersions {
		m.prowContext.OptionalPresubmit(pkg.E2ETest(m.prowContext, secondaryVersion, m.e2eCPURequest, m.e2eMemoryRequest))
	}

	m.prowContext.RequiredPresubmit(pkg.E2ETest(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest))

	m.prowContext.RequiredPresubmit(pkg.UpgradeTest(m.prowContext, m.primaryKubernetesVersion))

	m.prowContext.OptionalPresubmitIfChanged(pkg.LicenseTest(m.prowContext), `go.mod`)

	m.prowContext.OptionalPresubmit(pkg.E2ETestVenafiTPP(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest))
	m.prowContext.OptionalPresubmit(pkg.E2ETestVenafiCloud(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest))
	m.prowContext.OptionalPresubmit(pkg.E2ETestFeatureGatesDisabled(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest))
	m.prowContext.OptionalPresubmit(pkg.E2ETestWithBestPracticeInstall(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest))

	allKubernetesVersions := append(m.otherKubernetesVersions, m.primaryKubernetesVersion)

	if !m.isPreMakefileModules {
		m.prowContext.Periodics(pkg.MakeTest(m.prowContext), 2)
	} else {
		m.prowContext.Periodics(pkg.MakeTestOld(m.prowContext), 2)
	}

	// TODO: add chart periodic test?

	for _, kubernetesVersion := range allKubernetesVersions {
		m.prowContext.Periodics(pkg.E2ETest(m.prowContext, kubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest), 2)

	}

	m.prowContext.Periodics(pkg.E2ETestVenafiBoth(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest), 12)

	m.prowContext.Periodics(pkg.UpgradeTest(m.prowContext, m.primaryKubernetesVersion), 8)

	m.prowContext.Periodics(pkg.E2ETestWithBestPracticeInstall(m.prowContext, m.primaryKubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest), 24)

	for _, kubernetesVersion := range allKubernetesVersions {
		// TODO: roll this into above for loop; we have two for loops here to preserve the
		// ordering of the tests in the output file, making it easier to review the
		// differences between generated tests and existing handwritten tests
		m.prowContext.Periodics(pkg.E2ETestFeatureGatesDisabled(m.prowContext, kubernetesVersion, m.e2eCPURequest, m.e2eMemoryRequest), 24)
	}

	// Apply the default set of container names if none have been specified
	// TODO: this is the set that we want to migrate to in the future
	if m.containerNames == nil {
		m.containerNames = []string{"controller", "acmesolver", "startupapicheck", "cainjector", "webhook"}
	}

	for _, container := range m.containerNames {
		m.prowContext.Periodics(pkg.TrivyTest(m.prowContext, container), 24)
	}

	return m.prowContext.JobFile()
}

// KnownBranches returns a list of all branches which have been configured here
func KnownBranches() []string {
	var availableBranches []string

	for branch := range knownBranches {
		availableBranches = append(availableBranches, branch)
	}

	slices.Sort(availableBranches)
	return availableBranches
}

// SpecForBranch returns a spec for the named branch, if it exists
func SpecForBranch(originalBranch string) (BranchSpec, error) {
	branch := strings.ToLower(originalBranch)

	spec, ok := knownBranches[branch]
	if !ok {
		return BranchSpec{}, fmt.Errorf("unknown branch %q; known branches are %q", originalBranch, KnownBranches())
	}

	return spec, nil
}
