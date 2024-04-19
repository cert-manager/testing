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

// There are upstream definitions of these structs here:
// https://github.com/kubernetes/test-infra/blob/857418c31f6963014ac8821c63e1053c2c0e7e88/prow/config/jobs.go

// Rather than importing the prow struct definitions (and pulling in a bunch of dependencies)
// we copy the structs + fields we actually use in practice here

type JobFile struct {
	Presubmits map[string][]*PresubmitJob `yaml:"presubmits"`
	Periodics  []*PeriodicJob             `yaml:"periodics"`
}

type Job struct {
	Name string `yaml:"name"`

	MaxConcurrency int `yaml:"max_concurrency"`

	Decorate bool `yaml:"decorate"`

	Annotations map[string]string `yaml:"annotations"`

	Labels map[string]string `yaml:"labels"`

	Spec JobSpec `yaml:"spec"`
}

type JobSpec struct {
	Containers []Container `yaml:"containers"`
	DNSConfig  DNSConfig   `yaml:"dnsConfig"`
}

type Container struct {
	Image string `yaml:"image"`

	Args []string `yaml:"args"`

	Resources ContainerResources `yaml:"resources"`

	SecurityContext *SecurityContext `yaml:"securityContext,omitempty"`

	Lifecycle *Lifecycle `yaml:"lifecycle,omitempty"`
}

type ContainerResources struct {
	Requests ContainerResourceRequest `yaml:"requests"`
}

type ContainerResourceRequest struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type DNSConfig struct {
	Options []DNSConfigOption `yaml:"options"`
}

type DNSConfigOption struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func DefaultDNSConfig() DNSConfig {
	return DNSConfig{
		Options: []DNSConfigOption{
			{
				Name:  "ndots",
				Value: "1",
			},
		},
	}
}

type SecurityContext struct {
	Privileged bool `yaml:"privileged"`

	Capabilities *SecurityContextCapabilities `yaml:"capabilities,omitempty"`
}

type SecurityContextCapabilities struct {
	Add []string `yaml:"add"`
}

type Lifecycle struct {
	PreStop LifecycleHandler `yaml:"preStop"`
}

type LifecycleHandler struct {
	Exec ExecAction `yaml:"exec"`
}

type ExecAction struct {
	Command []string `yaml:"command"`
}

type PresubmitJob struct {
	Job `yaml:",inline"`

	Branches []string `yaml:"branches"`

	AlwaysRun bool `yaml:"always_run"`
	Optional  bool `yaml:"optional"`

	RunIfChanged string `yaml:"run_if_changed,omitempty"`
}

type PeriodicJob struct {
	Job `yaml:",inline"`

	ExtraRefs []ExtraRef `yaml:"extra_refs"`

	// PeriodicityHours isn't present in upstream structs; we use it to track when
	// a job should recur. Keeping track of this allows us to spread jobs across
	// time after we know how many jobs we have to avoid several being scheduled at once.
	PeriodicityHours int `yaml:"-"`

	Cron string `yaml:"cron,omitempty"`
}

type ExtraRef struct {
	Org  string `yaml:"org"`
	Repo string `yaml:"repo"`

	BaseRef string `yaml:"base_ref"`
}
