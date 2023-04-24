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
	"strconv"
	"time"
)

// ProwContext holds jobs and information required to configure jobs for a given release channel.
type ProwContext struct {
	// Branch is the name of the branch corresponding to the release channel modelled by this ProwContext.
	// While it's possible to define a presubmit for multiple branches, this often doesn't correctly model
	// how cert-manager uses prow in practice - usually, we want a different set of supported kubernetes
	// versions for each major cert-manager release (and therefore branch), and in any case want a different
	// dashboard for each supported release channel.
	Branch string

	// Image is the common test image used for running prow jobs.
	Image string

	// PresubmitDashboard, if set, will generate a presubmit TestGrid dashboard name based on the branch name
	// for each presubmit job. If false, no presubmits will be added to a TestGrid dashboard.
	PresubmitDashboard bool

	// PeriodicDashboard, if set, will generate a periodic TestGrid dashboard name based on the branch name
	// for each periodic job. If false, no periodics will be added to a TestGrid dashboard.
	PeriodicDashboard bool

	// Org is the GitHub organisation of the repository under test.
	Org string

	// Repo is the GitHub repository name of the repository under test.
	Repo string

	presubmits []*PresubmitJob
	periodics  []*PeriodicJob

	minutesCounter time.Time
}

// RequiredPresubmit adds a presubmit which is run by default and required to pass before a PR can be merged
func (pc *ProwContext) RequiredPresubmit(job *Job) {
	pc.addPresubmit(job, true, false, "")
}

// RequiredPresubmits adds a list of jobs to the context
func (pc *ProwContext) RequiredPresubmits(jobs []*Job) {
	for _, job := range jobs {
		pc.addPresubmit(job, true, false, "")
	}
}

// OptionalPresubmit adds a presubmit which is not run by default and is optional
func (pc *ProwContext) OptionalPresubmit(job *Job) {
	pc.addPresubmit(job, false, true, "")
}

// OptionalPresubmitIfChanged adds a presubmit which is not run by default and is optional unless a file has been
// changed which matches changedFileRegex. In that situation, the job is always run.
// See https://docs.prow.k8s.io/docs/jobs/#triggering-jobs-based-on-changes
func (pc *ProwContext) OptionalPresubmitIfChanged(job *Job, changedFileRegex string) {
	pc.addPresubmit(job, false, true, changedFileRegex)
}

func (pc *ProwContext) addPresubmit(job *Job, alwaysRun bool, optional bool, changedFileRegex string) {
	job.Name = pc.presubmitJobName(job.Name)

	if pc.PresubmitDashboard {
		addTestGridAnnotations(pc.presubmitDashboardName())(job)
	}

	pc.presubmits = append(pc.presubmits, &PresubmitJob{
		Job: *job,
		// see the comment on ProwContext.Branch for why we only support a single branch here
		Branches:     []string{pc.Branch},
		AlwaysRun:    alwaysRun,
		Optional:     optional,
		RunIfChanged: changedFileRegex,
	})
}

// Periodic adds periodic jobs which will run every `periodicityHours` hours, at some minute
// within the hour, one job for each configured branch
func (pc *ProwContext) Periodics(job *Job, periodicityHours int) {
	originalName := job.Name

	job.Name = pc.periodicJobName(originalName)

	if pc.PeriodicDashboard {
		addTestGridAnnotations(pc.periodicDashboardName())(job)
	}

	pc.periodics = append(pc.periodics, &PeriodicJob{
		Job: *job,
		ExtraRefs: []ExtraRef{
			{
				Org:     pc.Org,
				Repo:    pc.Repo,
				BaseRef: pc.Branch,
			},
		},
		Interval: strconv.Itoa(periodicityHours) + "h",
		// TODO: use Cron instead of Interval
		// Cron: pc.cronSchedule(periodicityHours),
	})
}

func (pc *ProwContext) JobFile() *JobFile {
	// TODO: when using Cron instead of Interval for periodics, adjust all periodics
	// here to spread them evenly throughout the hour

	presubmitKey := fmt.Sprintf("%s/%s", pc.Org, pc.Repo)

	return &JobFile{
		Presubmits: map[string][]*PresubmitJob{
			presubmitKey: pc.presubmits,
		},
		Periodics: pc.periodics,
	}
}

// presubmitJobName returns a prow name for the given presubmit job. For example,
// for the branch "release-1.0" and the test "foo", this would return "pull-cert-manager-release-1.0-foo"
func (pc *ProwContext) presubmitJobName(name string) string {
	return fmt.Sprintf("pull-%s-%s-%s", pc.Repo, pc.Branch, name)
}

// periodicJobName returns a prow name for the given periodic job. For example,
// for the branch "release-1.0" and the test "foo", this would return "ci-cert-manager-release-1.0-foo"
func (pc *ProwContext) periodicJobName(name string) string {
	return fmt.Sprintf("ci-%s-%s-%s", pc.Repo, pc.Branch, name)
}

func (pc *ProwContext) presubmitDashboardName() string {
	return fmt.Sprintf("%s-presubmits-%s", pc.Repo, pc.Branch)
}

func (pc *ProwContext) periodicDashboardName() string {
	return fmt.Sprintf("%s-periodics-%s", pc.Repo, pc.Branch)
}

func (pc *ProwContext) cronSchedule(periodicityHours int) string {
	minute := pc.minutesValue()

	return fmt.Sprintf("*/%d %d * * *", minute, periodicityHours)
}

// minutesValue returns a minute value (0 - 59) at which a test should be run and then
// increases the next value returned. This helps to prevent every test running at the same
// minute within the hour causing a spiky distribution of tests.
func (pc *ProwContext) minutesValue() int {
	minuteVal := pc.minutesCounter.Minute()

	pc.minutesCounter = pc.minutesCounter.Add(4 * time.Minute)

	return minuteVal
}
