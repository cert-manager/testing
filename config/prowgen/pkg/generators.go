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
	"strings"
)

// MakeTest generates a test which runs linting and verification targets as well as
// unit and integration tests
func MakeTest(ctx *ProwContext) *Job {
	job := jobTemplate(
		"make-test",
		"Runs unit and integration tests and verification scripts",
		addLocalCacheLabel,
		addGoCacheLabel,
		addMaxConcurrency(8),
	)

	makeJobs, cpuRequest := calculateMakeConcurrency("2000m")

	job.Spec.Containers = []Container{
		{
			Image: ctx.Image,
			Args: []string{
				"runner",
				"make",
				makeJobs,
				"vendor-go",
				"ci-presubmit",
				"test-ci",
			},
			Resources: ContainerResources{
				Requests: ContainerResourceRequest{
					CPU:    cpuRequest,
					Memory: "4Gi",
				},
			},
		},
	}

	return job
}

// ChartTest generates a test which lints helm charts. This is run inside a container
// and so requires additional permissions.
func ChartTest(ctx *ProwContext) *Job {
	job := jobTemplate(
		"chart",
		"Verifies the Helm chart passes linting checks",
		addDindLabel,
		addLocalCacheLabel,
		addGoCacheLabel,
		addMaxConcurrency(8),
	)

	job.Spec.Containers = []Container{
		{
			Image: ctx.Image,
			Args: []string{
				"runner",
				"make",
				"vendor-go",
				"verify-chart",
			},
			Resources: ContainerResources{
				Requests: ContainerResourceRequest{
					CPU:    "1",
					Memory: "1Gi",
				},
			},
			SecurityContext: &SecurityContext{
				Privileged: true,
			},
		},
	}

	return job
}

// LicenseTest generates a test which validates the LICENSES file. Since the verify-licenses make target
// depends on external services for license checking, running it on every PR would introduce the possibilities
// for flakes if, say, a vanity import site such as gopkg.in was down.
// We special case the license test so it only runs when go.mod has changed.
func LicenseTest(ctx *ProwContext) *Job {
	job := jobTemplate(
		"license",
		"Verifies LICENSES are up to date; only needs to be run if go.mod has changed",
		addLocalCacheLabel,
		addGoCacheLabel,
		addMaxConcurrency(8),
	)

	job.Spec.Containers = []Container{
		{
			Image: ctx.Image,
			Args: []string{
				"runner",
				"make",
				"vendor-go",
				"verify-licenses",
			},
			Resources: ContainerResources{
				Requests: ContainerResourceRequest{
					CPU:    "1",
					Memory: "1Gi",
				},
			},
		},
	}

	return job
}

// E2ETest generates a test which runs end-to-end tests with feature gates enabled. This
// is run inside a container and requires additional permissions.
func E2ETest(ctx *ProwContext, k8sVersion string, cpuRequest, memoryRequest string) *Job {
	// we don't want to use dots in names, so replace with dashes
	nameVersion := strings.ReplaceAll(k8sVersion, ".", "-")

	desc := fmt.Sprintf("Runs the end-to-end test suite against a Kubernetes v%s cluster", k8sVersion)

	job := jobTemplate(
		"e2e-v"+nameVersion,
		desc,
		addDindLabel,
		addCloudflareCredentialsLabel,
		addLocalCacheLabel,
		addGoCacheLabel,
		addStandardE2ELabels(k8sVersion),
		addRetryFlakesLabel,
		addMaxConcurrency(4),
	)

	makeJobs, cpuRequest := calculateMakeConcurrency(cpuRequest)

	k8sVersionArg := fmt.Sprintf("K8S_VERSION=%s", k8sVersion)

	job.Spec.Containers = []Container{
		{
			Image: ctx.Image,
			Args: []string{
				"runner",
				"make",
				makeJobs,
				"vendor-go",
				"e2e-ci",
				k8sVersionArg,
			},
			Resources: ContainerResources{
				Requests: ContainerResourceRequest{
					CPU:    cpuRequest,
					Memory: memoryRequest,
				},
			},
			SecurityContext: &SecurityContext{
				Privileged: true,
				Capabilities: &SecurityContextCapabilities{
					Add: []string{"SYS_ADMIN"},
				},
			},
			Lifecycle: &Lifecycle{
				PreStop: LifecycleHandler{
					Exec: ExecAction{
						Command: []string{
							"/bin/sh",
							"-c",
							"make kind-logs",
						},
					},
				},
			},
		},
	}

	return job
}

// E2ETestVenafiTPP generates a test which runs end-to-end tests only focusing on Venafi TPP.
// This runs inside a container and so requires additional permissions.
func E2ETestVenafiTPP(ctx *ProwContext, k8sVersion string, cpuRequest, memoryRequest string) *Job {
	job := E2ETest(ctx, k8sVersion, cpuRequest, memoryRequest)

	job.Name = job.Name + "-issuers-venafi-tpp"
	job.Annotations["description"] = "Runs the E2E tests with 'Venafi TPP' in name"

	job.Labels = make(map[string]string)

	addDindLabel(job)
	addLocalCacheLabel(job)
	addGoCacheLabel(job)
	addRetryFlakesLabel(job)
	addVenafiTPPLabels(job)

	return job
}

// E2ETestVenafiCloud generates a test which runs end-to-end tests only focusing on Venafi Cloud.
// This runs inside a container and so requires additional permissions.
func E2ETestVenafiCloud(ctx *ProwContext, k8sVersion string, cpuRequest, memoryRequest string) *Job {
	job := E2ETest(ctx, k8sVersion, cpuRequest, memoryRequest)

	job.Name = job.Name + "-issuers-venafi-cloud"
	job.Annotations["description"] = "Runs the E2E tests with 'Venafi Cloud' in name"

	job.Labels = make(map[string]string)

	addDindLabel(job)
	addLocalCacheLabel(job)
	addGoCacheLabel(job)
	addRetryFlakesLabel(job)
	addVenafiCloudLabels(job)

	return job
}

// E2ETestVenafiBoth generates a test which runs end-to-end tests focusing on
// both Venafi TPP and Venafi Cloud.
// This runs inside a container and so requires additional permissions.
func E2ETestVenafiBoth(ctx *ProwContext, k8sVersion string, cpuRequest, memoryRequest string) *Job {
	job := E2ETest(ctx, k8sVersion, cpuRequest, memoryRequest)

	job.Name = job.Name + "-issuers-venafi"
	job.Annotations["description"] = "Runs Venafi (VaaS and TPP) e2e tests"

	job.Labels = make(map[string]string)

	addDindLabel(job)
	addLocalCacheLabel(job)
	addGoCacheLabel(job)
	addRetryFlakesLabel(job)
	addVenafiBothLabels(job)

	return job
}

// E2ETestFeatureGatesDisabled generates a test which runs e2e tests with feature gates disabled
func E2ETestFeatureGatesDisabled(ctx *ProwContext, k8sVersion string, cpuRequest, memoryRequest string) *Job {
	job := E2ETest(ctx, k8sVersion, cpuRequest, memoryRequest)

	job.Name = job.Name + "-feature-gates-disabled"
	job.Annotations["description"] = "Runs the E2E tests with all feature gates disabled"

	job.Labels = make(map[string]string)

	addCloudflareCredentialsLabel(job)
	addDindLabel(job)
	addDisableFeatureGatesLabel(job)
	addGinkgoSkipDefaultLabel(job)
	addLocalCacheLabel(job)
	addGoCacheLabel(job)
	addRetryFlakesLabel(job)

	return job
}

// E2ETestWithBestPracticeInstall generates a test which runs e2e tests
// with cert-manager installed in accordance with
// https://cert-manager.io/docs/installation/best-practice/
func E2ETestWithBestPracticeInstall(ctx *ProwContext, k8sVersion string, cpuRequest, memoryRequest string) *Job {
	job := E2ETest(ctx, k8sVersion, cpuRequest, memoryRequest)

	job.Name = job.Name + "-bestpractice-install"
	job.Annotations["description"] = "Runs the E2E tests with cert-manager installed in accordance with https://cert-manager.io/docs/installation/best-practice/"

	job.Labels = make(map[string]string)

	addCloudflareCredentialsLabel(job)
	addDindLabel(job)
	addDisableFeatureGatesLabel(job)
	addGinkgoSkipDefaultLabel(job)
	addLocalCacheLabel(job)
	addGoCacheLabel(job)
	addRetryFlakesLabel(job)
	addBestPracticeInstallLabel(job)

	return job
}

// UpgradeTest generates a test which tests an upgrade from the latest released version
// of cert-manager to the version specified by the test ref / branch. This test runs
// inside a container and so requires additional privileges.
func UpgradeTest(ctx *ProwContext, k8sVersion string) *Job {
	nameVersion := strings.ReplaceAll(k8sVersion, ".", "-")

	job := jobTemplate(
		"e2e-v"+nameVersion+"-upgrade",
		"Runs cert-manager upgrade from latest published release",
		addDindLabel,
		addLocalCacheLabel,
		addGoCacheLabel,
		addMaxConcurrency(4),
	)

	k8sVersionArg := fmt.Sprintf("K8S_VERSION=%s", k8sVersion)

	job.Spec.Containers = []Container{
		{
			Image: ctx.Image,
			Args: []string{
				"runner",
				"make",
				k8sVersionArg,
				"vendor-go",
				"test-upgrade",
			},
			Resources: ContainerResources{
				Requests: ContainerResourceRequest{
					CPU:    "3500m",
					Memory: "6Gi",
				},
			},
			SecurityContext: &SecurityContext{
				Privileged: true,
				Capabilities: &SecurityContextCapabilities{
					Add: []string{"SYS_ADMIN"},
				},
			},
		},
	}

	return job
}

// TrivyTest generates a test which runs a Trivy scan of a built container image which matches the given name.
// Note that there's also a "make trivy-scan-all" target, but this will fail as soon as one of its dependencies fails,
// so e.g. if there's a vuln in the "controller" container we might never scan "ctl" container.
// Instead, we generate a test for each container so it's obvious which ones have failures and it's easier to get results
// for each container
func TrivyTest(ctx *ProwContext, containerName string) *Job {
	containerName = strings.ToLower(containerName)

	job := jobTemplate(
		fmt.Sprintf("trivy-test-%s", containerName),
		fmt.Sprintf("Runs a Trivy scan against the %s container", containerName),
		addLocalCacheLabel,
		addGoCacheLabel,
		addDindLabel,
		addMaxConcurrency(2),
		// Need to ensure that trivy tests send a failure email as soon as they fail since
		// they tend to be run relatively infrequently and a failure is important to address
		addTestGridCustomFailuresToAlert(1),
		// Ask TestGrid to alert us if the job hasn't run in the last 36 hours. Sets
		// an upper limit on how regularly the job can be scheduled.
		addTestGridStaleResultsAlert(36),
	)

	makeJobs, cpuRequest := calculateMakeConcurrency("1000m")

	job.Spec.Containers = []Container{
		{
			Image: ctx.Image,
			Args: []string{
				"runner",
				"make",
				makeJobs,
				"vendor-go",
				fmt.Sprintf("trivy-scan-%s", containerName),
			},
			Resources: ContainerResources{
				Requests: ContainerResourceRequest{
					CPU:    cpuRequest,
					Memory: "2Gi",
				},
			},
			SecurityContext: &SecurityContext{
				Privileged: true,
			},
		},
	}

	return job
}
