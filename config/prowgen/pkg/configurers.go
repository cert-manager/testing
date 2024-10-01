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

import "fmt"

type JobConfigurer func(*Job)

// jobTemplate defines a 'default' job, where standard parameters can be set. All jobs
// should have a name and a friendly description of what they do.
// Callers can also pass a list of "configurers" which will modify the template before
// it's returned for use.
func jobTemplate(name string, description string, configurers ...JobConfigurer) *Job {
	job := &Job{
		Name:     name,
		Decorate: true,
		Annotations: map[string]string{
			"description": description,
		},
		Labels: make(map[string]string),
		Spec: JobSpec{
			// Fully omit the kube-dns DNS server, and directly talk to
			// Google's public DNS server. Generally, the tests don't need
			// to access any in-cluster services, so not being able to resolve
			// xxxxx.cluster.local queries is ok.
			DNSPolicy: "None",
			DNSConfig: DNSConfig{
				Nameservers: []string{"8.8.8.8", "8.8.4.4"},
			},
		},
	}

	for _, c := range configurers {
		c(job)
	}

	return job
}

func addTrivyLabel(job *Job) {
	job.Labels["preset-trivy"] = "true"
}

func addLocalCacheLabel(job *Job) {
	job.Labels["preset-local-cache"] = "true"
}

func addGoCacheLabel(job *Job) {
	job.Labels["preset-go-cache"] = "true"
}

func addDindLabel(job *Job) {
	job.Labels["preset-dind-enabled"] = "true"
}

func addCloudflareCredentialsLabel(job *Job) {
	job.Labels["preset-cloudflare-credentials"] = "true"
}

func addRetryFlakesLabel(job *Job) {
	job.Labels["preset-retry-flakey-jobs"] = "true"
}

func addGinkgoSkipDefaultLabel(job *Job) {
	job.Labels["preset-ginkgo-skip-default"] = "true"
}

func addDisableFeatureGatesLabel(job *Job) {
	job.Labels["preset-disable-all-alpha-beta-feature-gates"] = "true"
}

func addVenafiTPPLabels(job *Job) {
	job.Labels["preset-ginkgo-focus-venafi-tpp"] = "true"
	job.Labels["preset-venafi-tpp-credentials"] = "true"
}

func addVenafiBothLabels(job *Job) {
	job.Labels["preset-ginkgo-focus-venafi"] = "true"

	job.Labels["preset-venafi-cloud-credentials"] = "true"
	job.Labels["preset-venafi-tpp-credentials"] = "true"
}

func addVenafiCloudLabels(job *Job) {
	job.Labels["preset-ginkgo-focus-venafi-cloud"] = "true"
	job.Labels["preset-venafi-cloud-credentials"] = "true"
}

func addBestPracticeInstallLabel(job *Job) {
	job.Labels["preset-bestpractice-install"] = "true"
}

func addStandardE2ELabels(kubernetesVersion string) JobConfigurer {
	return func(job *Job) {
		addGinkgoSkipDefaultLabel(job)

		majorVersion, minorVersion, err := splitKubernetesVersion(kubernetesVersion)
		if err != nil {
			// note: we panic here because this tool is developer-facing and because an
			// error here suggests programmer error (e.g. a typo'd k8s version)
			// adding 'return nil' for every configurer - most of which shouldn't fail in
			// any reasonable scenario - seems far messier than a panic here
			panic(err)
		}

		if majorVersion == 1 && minorVersion < 22 {
			// SSA (server-side apply) is only fully supported in k8s 1.22+
			job.Labels["preset-enable-all-feature-gates-disable-ssa"] = "true"
			return
		}

		job.Labels["preset-enable-all-feature-gates"] = "true"
	}
}

// addTestGridAnnotations inserts standard testgrid annotations for the job.
// For a list of testgrid annotations, see:
// https://github.com/GoogleCloudPlatform/testgrid/blob/444774c4b660dad5ab3c1f47e0579d37deb6b5b0/config.md#prow-job-configuration
func addTestGridAnnotations(dashboardName string) JobConfigurer {
	return func(job *Job) {
		job.Annotations["testgrid-create-job-group"] = "true"
		job.Annotations["testgrid-dashboards"] = dashboardName
		job.Annotations["testgrid-alert-email"] = AlertEmailAddress
	}
}

// addTestGridCustomFailuresToAlert changes the number of failures required before TestGrid
// marks a job as "failed" (rather thank "flaky")
func addTestGridCustomFailuresToAlert(failuresToAlert int) JobConfigurer {
	return func(job *Job) {
		job.Annotations["testgrid-num-failures-to-alert"] = fmt.Sprintf("%d", failuresToAlert)
	}
}

// addTestGridStaleResultsAlert sets, in hours, the length of time before a job should be
// considered stale. This guards against a job not running for whatever reason.
func addTestGridStaleResultsAlert(hoursUntilStale int) JobConfigurer {
	return func(job *Job) {
		job.Annotations["testgrid-alert-stale-results-hours"] = fmt.Sprintf("%d", hoursUntilStale)
	}
}

func addMaxConcurrency(maxConcurrency int) JobConfigurer {
	return func(job *Job) {
		job.MaxConcurrency = maxConcurrency
	}
}
